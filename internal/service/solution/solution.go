package solution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"xinde/internal/dao/account"
	"xinde/internal/dao/attachment"
	"xinde/internal/dao/device"
	"xinde/internal/dao/solution"
	deviceDto "xinde/internal/dto/device"
	dto "xinde/internal/dto/solution"
	attachmentModel "xinde/internal/model/attachment"
	deviceModel "xinde/internal/model/device"
	"xinde/pkg/jwt"
)

type Service struct {
	dao           *solution.Dao
	deviceDao     *device.Dao
	attachmentDao *attachment.Dao
	j             *jwt.JWTService
	accountDao    *account.Dao
}

func NewSolutionService() (*Service, error) {
	deviceDao, err := device.NewDeviceDao()
	if err != nil {
		return nil, fmt.Errorf("NewDeviceDao() 创建Dao实例失败: %v", err)
	}
	dao, err := solution.NewSolutionDao()
	if err != nil {
		return nil, fmt.Errorf("NewSolutionDao() 创建Dao实例失败: %v", err)
	}
	attachmentDao, err := attachment.NewAttachmentDao()
	if err != nil {
		return nil, fmt.Errorf("NewAttachmentDao() 创建Dao实例失败: %v", err)
	}
	accountDao, err := account.NewRegisterDao()
	if err != nil {
		return nil, fmt.Errorf("NewRegisterDao() 创建Dao实例失败: %v", err)
	}
	j := jwt.NewJWTService()
	return &Service{
		dao:           dao,
		deviceDao:     deviceDao,
		attachmentDao: attachmentDao,
		j:             j,
		accountDao:    accountDao,
	}, nil
}

func (s *Service) Query(userID uint, req *dto.QueryReq) (*dto.QueryResp, error) {
	// 1. 查询方案列表和总数
	total, solutions, err := s.dao.QuerySolutions(s.dao.DB(), req)
	if err != nil {
		return nil, err
	}

	// 2. 聚合可用筛选条件
	aggFilters, err := s.dao.AggregateFilters(s.dao.DB(), req)
	if err != nil {
		return nil, err
	}

	// 3. 聚合外部数据 (价格 & API)
	solutionDataList, err := s.aggregateExternalData(userID, solutions)
	if err != nil {
		return nil, err
	}

	// 4. 组装可用筛选条件
	availableFilters, err := s.buildAvailableFilters(req.DeviceTypeID, aggFilters)
	if err != nil {
		return nil, err
	}

	// 5. 组装最终响应
	resp := &dto.QueryResp{
		Solutions: &dto.SolutionsPageData{
			List:     solutionDataList,
			Total:    total,
			Page:     req.Pagination.Page,
			PageSize: req.Pagination.PageSize,
			Pages:    (total + int64(req.Pagination.PageSize) - 1) / int64(req.Pagination.PageSize),
		},
		AvailableFilters: availableFilters,
	}

	return resp, nil
}

// --- 新增：用于调用 API 的辅助结构体 ---

type ApiRequestDetail struct {
	Limitelength string `json:"limitelength"`
}

type ApiRequestBody struct {
	Dbname  string           `json:"dbname"`
	Queryid string           `json:"queryid"`
	Detail  ApiRequestDetail `json:"detail"`
}

type ApiResultData struct {
	Brand    string      `json:"brand"`
	Bsonhand float64     `json:"bsonhand"`
	Itemcode string      `json:"itemcode"`
	Onhand   interface{} `json:"onhand"` // 使用 interface{} 来处理 null 或数字
	Pic      string      `json:"pic"`
}

type ApiResponse struct {
	Data   []ApiResultData `json:"data"`
	Errmsg string          `json:"errmsg"`
	Errno  string          `json:"errno"`
}

// aggregateExternalData 是新的辅助函数，负责将 model 转换为包含外部数据的 DTO
func (s *Service) aggregateExternalData(userID uint, solutions []*deviceModel.Device) ([]*dto.SolutionData, error) {
	var solutionDataList []*dto.SolutionData

	// 1. 收集所有不重复的 product_code
	productCodeSet := make(map[string]bool)
	for _, sol := range solutions {
		var importDetails deviceDto.ImportDetailsDTO
		if json.Unmarshal(sol.Details, &importDetails) == nil {
			for _, comp := range importDetails.Components {
				if comp.ProductCode != "" {
					productCodeSet[comp.ProductCode] = true
				}
			}
		}
	}
	var productCodes []string
	for code := range productCodeSet {
		productCodes = append(productCodes, code)
	}

	// 2. 批量调用二方 API
	apiDataMap, err := s.callExternalAPI(productCodes)
	if err != nil {
		// 如果 API 调用失败，我们可以选择返回错误，或者只记录日志并继续（返回不带API数据的结果）
		// 这里我们选择返回错误
		return nil, fmt.Errorf("调用二方服务失败: %w", err)
	}

	// [占位符] 3. 批量查询 MySQL 价格表
	priceResults, err := s.accountDao.FindPricesForUser(s.accountDao.DB(), userID, productCodes)
	if err != nil {
		return nil, fmt.Errorf("查询价格失败: %w", err)
	}

	// 将价格结果转换为 product_code -> price 的 map，方便查找
	priceMap := make(map[string]float64)
	for _, p := range priceResults {
		switch p.PriceLevel {
		case "price_1":
			priceMap[p.ProductCode] = p.Price1
		case "price_2":
			priceMap[p.ProductCode] = p.Price2
		case "price_3":
			priceMap[p.ProductCode] = p.Price3
		case "price_4":
			priceMap[p.ProductCode] = p.Price4
		default:
			priceMap[p.ProductCode] = p.Price1 // 默认价格
		}
	}

	// 4. 遍历并聚合数据
	for _, sol := range solutions {
		var importDetails deviceDto.ImportDetailsDTO
		if err := json.Unmarshal(sol.Details, &importDetails); err != nil {
			continue
		}

		readDetails := &dto.DetailsData{
			Filters:    importDetails.Filters,
			Parameters: importDetails.Parameters,
			Components: []*dto.ComponentData{},
		}

		for _, comp := range importDetails.Components {
			readComp := &dto.ComponentData{
				Name:        comp.Name,
				ProductCode: comp.ProductCode,
				SpecCode:    comp.SpecCode,
				Price:       priceMap[comp.ProductCode], // 从价格 map 中获取
			}

			// 从 API 结果中填充数据
			if apiData, ok := apiDataMap[comp.ProductCode]; ok {
				readComp.Brand = apiData.Brand

				// 处理 onhand (可能为 null)
				if onhandVal, ok := apiData.Onhand.(float64); ok {
					readComp.InventoryXinde = fmt.Sprintf("%.0f", onhandVal)
				} else {
					readComp.InventoryXinde = "0" // 或者 "0"
				}

				readComp.InventoryGongpin = fmt.Sprintf("%.0f", apiData.Bsonhand)

				// 拼接图片 URL
				if apiData.Pic != "" {
					imageBaseURL := viper.GetString("external_api.image_base_url")
					readComp.ImageURL = imageBaseURL + strings.TrimPrefix(apiData.Pic, "/")
				}
			}
			readDetails.Components = append(readDetails.Components, readComp)
		}

		solutionDataList = append(solutionDataList, &dto.SolutionData{
			ID:      sol.ID,
			Name:    sol.Name,
			Details: readDetails,
		})
	}

	return solutionDataList, nil
}

// callExternalAPI 是一个私有方法，用于调用二方服务
func (s *Service) callExternalAPI(productCodes []string) (map[string]ApiResultData, error) {
	if len(productCodes) == 0 {
		return make(map[string]ApiResultData), nil
	}

	// 1. 准备请求体
	reqBody := ApiRequestBody{
		Dbname:  viper.GetString("external_api.dbname"),
		Queryid: viper.GetString("external_api.queryid"),
		Detail: ApiRequestDetail{
			Limitelength: strings.Join(productCodes, ","),
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化API请求体失败: %w", err)
	}

	// 2. 创建 HTTP 请求
	apiURL := viper.GetString("external_api.url")
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 4. 读取和解析响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取API响应体失败: %w", err)
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("反序列化API响应失败: %w", err)
	}

	// 5. 检查业务错误
	if apiResp.Errno != "0" {
		return nil, fmt.Errorf("API返回业务错误: %s", apiResp.Errmsg)
	}

	// 6. 将结果转换为 map 方便查找
	resultMap := make(map[string]ApiResultData)
	for _, item := range apiResp.Data {
		resultMap[item.Itemcode] = item
	}

	return resultMap, nil
}

// 【新增实现】buildAvailableFilters
// buildAvailableFilters builds the filter structure with associated images.
func (s *Service) buildAvailableFilters(deviceTypeID uint, aggFilters map[string][]string) ([]*dto.AvailableFilter, error) {
	var availableFilters []*dto.AvailableFilter

	// 1. 从 PG 获取此设备类型下所有的图片配置
	filterImages, err := s.deviceDao.GetFilterImagesByDeviceTypeID(s.deviceDao.DB(), deviceTypeID)
	if err != nil {
		return nil, fmt.Errorf("获取筛选条件图片失败: %w", err)
	}

	// 2. 将图片配置转换为 value -> url 的 map，方便快速查找
	filterImageMap := make(map[string]string)
	if len(filterImages) > 0 {
		// a. 收集所有 filter_image 的主键 ID
		var filterImageIDs []uint
		for _, fi := range filterImages { // 【已修复】现在 fi 被使用了
			filterImageIDs = append(filterImageIDs, fi.ID)
		}

		// b. 从 MySQL 批量查询这些 ID 对应的附件
		businessType := viper.GetString("business_type.filter_image")
		// 【注意】attachmentDao 使用的是 MySQL 的 DB 实例
		attachments, err := s.attachmentDao.GetAttachmentsByBusinessAndIDs(s.attachmentDao.DB(), businessType, filterImageIDs)
		if err != nil {
			return nil, fmt.Errorf("从MySQL获取附件信息失败: %w", err)
		}

		// c. 创建一个 business_id (即 filter_image.id) -> attachment 的映射
		attachmentMap := make(map[uint]*attachmentModel.Attachment)
		for _, att := range attachments {
			attachmentMap[att.BusinessID] = att
		}

		// d. 构建最终的 filter_value -> URL 映射
		baseURL := viper.GetString("server.base_url")
		urlPrefix := viper.GetString("attachment.upload_url_prefix")
		for _, fi := range filterImages {
			if att, ok := attachmentMap[fi.ID]; ok {
				filterImageMap[fi.FilterValue] = fmt.Sprintf("%s%s/%s", baseURL, urlPrefix, att.StoragePath)
			}
		}
	}

	// 3. 遍历聚合出的筛选条件，组装最终结果
	for name, options := range aggFilters {
		filter := &dto.AvailableFilter{FilterName: name}

		// --- 【核心变更】对 options 进行排序 ---
		sort.SliceStable(options, func(i, j int) bool {
			// 尝试将选项作为数字进行比较
			numI, errI := strconv.ParseFloat(options[i], 64)
			numJ, errJ := strconv.ParseFloat(options[j], 64)

			// 如果两个都能成功转换为数字，则按数字大小排序
			if errI == nil && errJ == nil {
				return numI < numJ
			}

			// 否则，按标准的字符串字典序排序
			return options[i] < options[j]
		})
		// --- 排序结束 ---

		for _, optVal := range options {
			option := dto.FilterOption{Value: optVal}
			if url, ok := filterImageMap[optVal]; ok {
				option.ImageURL = url
			}
			filter.Options = append(filter.Options, option)
		}

		availableFilters = append(availableFilters, filter)
	}

	// 【可选】你也可以对最外层的 filter (按名称) 进行排序
	sort.Slice(availableFilters, func(i, j int) bool {
		return availableFilters[i].FilterName < availableFilters[j].FilterName
	})

	return availableFilters, nil
}
