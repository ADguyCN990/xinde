package solution

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"strings"
	"xinde/internal/dao/common"
	dto "xinde/internal/dto/solution"
	"xinde/internal/model/device"
	"xinde/internal/store"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func (d *Dao) DB() *gorm.DB {
	return d.db
}

func NewSolutionDao() (*Dao, error) {
	db := store.GetPDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化，请先调用 store.InitDB()")
	}
	commonDao, err := common.NewCommonPostgresDao()
	if err != nil {
		return nil, err
	}
	return &Dao{
		db:        db,
		commonDao: commonDao,
	}, nil
}

// buildDynamicQuery 根据req里的【已选定的filter】，动态构建where查询语句
func (d *Dao) buildDynamicQuery(tx *gorm.DB, req *dto.QueryReq) *gorm.DB {
	query := tx.Model(&device.Device{}).Where("device_type_id = ?", req.DeviceTypeID)

	// 【新增】创建一个 map 来暂存 min 和 max 的值，以便成对处理
	rangeFilters := make(map[string]map[string]interface{})

	for key, value := range req.CurrentFilters {
		// --- 【核心变更】处理范围筛选 ---
		if strings.HasSuffix(key, "_min") {
			baseName := strings.TrimSuffix(key, "_min")
			if _, ok := rangeFilters[baseName]; !ok {
				rangeFilters[baseName] = make(map[string]interface{})
			}
			rangeFilters[baseName]["min"] = value
			continue // 处理完后跳过，不进入下面的 switch
		}
		if strings.HasSuffix(key, "_max") {
			baseName := strings.TrimSuffix(key, "_max")
			if _, ok := rangeFilters[baseName]; !ok {
				rangeFilters[baseName] = make(map[string]interface{})
			}
			rangeFilters[baseName]["max"] = value
			continue // 处理完后跳过
		}

		// --- 处理精确匹配 (蓝色) ---
		switch v := value.(type) {
		case string:
			query = query.Where("details -> 'filters' ->> ? = ?", key, v)
		case float64:
			query = query.Where("details -> 'filters' ->> ? = ?", key, fmt.Sprintf("%g", v))
		// 可以添加对 int, bool 等类型的处理
		default:
			// 对于其他类型，保守地转为字符串进行比较
			query = query.Where("details -> 'filters' ->> ? = ?", key, fmt.Sprintf("%v", v))
		}
	}

	// --- 【新增】统一处理所有范围筛选 ---
	for baseName, minMax := range rangeFilters {
		minKey := baseName + "_min"
		maxKey := baseName + "_max"

		// 用户可能只提供了 min 或 max 之一
		if minVal, ok := minMax["min"]; ok {
			// (数据库中的 max 值) >= (用户输入的 min 值)
			query = query.Where("(details -> 'filters' ->> ?)::numeric >= ?", maxKey, minVal)
		}
		if maxVal, ok := minMax["max"]; ok {
			// (数据库中的 min 值) <= (用户输入的 max 值)
			query = query.Where("(details -> 'filters' ->> ?)::numeric <= ?", minKey, maxVal)
		}
	}

	return query
}

// QuerySolutions retrieves a paginated list of solutions based on dynamic filters.
func (d *Dao) QuerySolutions(tx *gorm.DB, req *dto.QueryReq) (int64, []*device.Device, error) {
	var total int64
	var solutions []*device.Device

	query := d.buildDynamicQuery(tx, req)

	// 先计算总数
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 再获取分页数据
	offset := (req.Pagination.Page - 1) * req.Pagination.PageSize
	err := query.Limit(req.Pagination.PageSize).Offset(offset).Find(&solutions).Error

	return total, solutions, err
}

// AggregateFilters aggregates the available filters from the result set.
func (d *Dao) AggregateFilters(tx *gorm.DB, req *dto.QueryReq) (map[string][]string, error) {

	// 获取所有符合当前筛选条件的item
	query := d.buildDynamicQuery(tx, req)

	// 【修正】使用 datatypes.JSON 来接收原始的 JSON 数据
	var results []struct {
		Filters datatypes.JSON `gorm:"column:filters"`
	}
	if err := query.Select("details -> 'filters' as filters").Scan(&results).Error; err != nil {
		return nil, err
	}

	/**
	*出来的结构类似这样:
	*{filters:{
	*'U钻类型' : '880型',
	*'工件材质': '1.钢件P'
	* ……
	*}}
	**/

	// 检查当前所有符合筛选条件的“方案”，统计出在这些方案中，还剩下哪些筛选条件以及它们有哪些可选的值
	//假设用户已经选择 U钻类型="880型"。现在数据库里还剩下 5 条符合这个条件的“方案”。
	//方案1 的 filters: {"U钻类型":"880型", "品牌":"博世", "钻孔深度_min":"10"}
	//方案2 的 filters: {"U钻类型":"880型", "品牌":"MA/美研", "钻孔深度_min":"12"}
	//方案3 的 filters: {"U钻类型":"880型", "品牌":"博世", "钻孔深度_min":"15"}
	//方案4 的 filters: {"U钻类型":"880型", "品牌":"MA/美研", "钻孔深度_min":"10"}
	//方案5 的 filters: {"U钻类型":"880型", "品牌":"山特维克", "钻孔深度_min":"12"}
	//处理方案1: aggMap 变成 {"U钻类型":{"880型":true}, "品牌":{"博世":true}, "钻孔深度_min":{"10":true}}
	//处理方案2: aggMap 变成 {"U钻类型":{"880型":true}, "品牌":{"博世":true, "MA/美研":true}, "钻孔深度_min":{"10":true, "12":true}}
	//处理方案3: aggMap 变成 {"U钻类型":{"880型":true}, "品牌":{"博世":true, "MA/美研":true}, "钻孔深度_min":{"10":true, "12":true, "15":true}}
	//处理方案4: 往 aggMap 里放 "MA/美研" 和 "10"，但因为它们已经存在了，所以 map 没变化。
	//处理方案5: aggMap 变成 {"U钻类型":{"880型":true}, "品牌":{"博世":true, "MA/美研":true, "山特维克":true}, "钻孔深度_min":{"10":true, "12":true, "15":true}}
	aggMap := make(map[string]map[string]bool)
	for _, res := range results {
		// 【新增】在 Go 代码中手动 Unmarshal
		var currentFilters map[string]interface{}
		if err := json.Unmarshal(res.Filters, &currentFilters); err != nil {
			// 记录日志并跳过无法解析的行
			// logger.Warn("无法解析 filters 字段: " + err.Error())
			continue
		}

		for key, val := range currentFilters {
			if _, ok := aggMap[key]; !ok {
				aggMap[key] = make(map[string]bool)
			}
			aggMap[key][fmt.Sprintf("%v", val)] = true
		}
	}

	// 转换格式
	finalMap := make(map[string][]string)
	for key, valSet := range aggMap {
		for val := range valSet {
			finalMap[key] = append(finalMap[key], val)
		}
	}

	return finalMap, nil
}
