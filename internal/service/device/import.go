package device

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
	"xinde/internal/dao/attachment"
	"xinde/internal/dao/device"
	dto "xinde/internal/dto/device"
	model "xinde/internal/model/attachment"
	deviceModel "xinde/internal/model/device"
	"xinde/pkg/jwt"
	"xinde/pkg/util"
)

type Service struct {
	dao           *device.Dao
	j             *jwt.JWTService
	attachmentDao *attachment.Dao
}

func NewDeviceService() (*Service, error) {
	dao, err := device.NewDeviceDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	attachmentDao, err := attachment.NewAttachmentDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	j := jwt.NewJWTService()
	return &Service{
		dao:           dao,
		j:             j,
		attachmentDao: attachmentDao,
	}, nil
}

func (s *Service) ImportFromExcel(adminID, groupID uint, deviceTypeName string, file, image *multipart.FileHeader) error {

	// 1. 解析Excel。这一步只做纯粹的解析，不涉及任何数据库或API调用。
	parsedData, err := s.parseFromExcel(file)
	if err != nil {
		return err
	}
	if len(parsedData) == 0 {
		return fmt.Errorf("excel没有解析到有效内容")
	}
	var deviceType *deviceModel.DeviceType

	// --- 2. 开启 PostgresSQL 事务，执行替换操作 ---
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {

		// a. 查找或创建 DeviceType
		deviceType, err = s.dao.FindOrCreateDeviceType(tx, deviceTypeName, groupID)
		if err != nil {
			return err
		}

		// b. 删除与此 DeviceType 关联的所有旧方案 (Device)
		err = s.dao.DeleteByDeviceTypeID(tx, deviceType.ID)
		if err != nil {
			return err
		}

		// c. 准备批量创建的新 "方案" (Device) 数据
		var solutions []*deviceModel.Device
		for i, data := range parsedData {
			deviceName := fmt.Sprintf("方案%d", i+1)
			detailJson, err := json.Marshal(data.Details)
			if err != nil {
				return fmt.Errorf("序列化方案的Detail失败: " + err.Error())
			}
			solution := &deviceModel.Device{
				Name:         deviceName,
				DeviceTypeID: deviceType.ID,
				Details:      detailJson,
			}
			solutions = append(solutions, solution)
		}

		// d. 批量写入新的 "方案" (Device) 到数据库
		if len(solutions) > 0 {
			err := s.dao.BatchCreateDevice(tx, solutions)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("导入设备提交事务失败: " + err.Error())
	}

	// --- 3. 【独立】处理主图附件 ---
	err = s.attachmentDao.DB().Transaction(func(tx *gorm.DB) error {

		// a. 获取business_type
		fileBusinessType := viper.GetString("business_type.device_import")
		iconBusinessType := viper.GetString("business_type.device_icon")

		// b. 查找并删除旧的数据库记录
		err := s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(tx, fileBusinessType, []uint{deviceType.ID})
		if err != nil {
			return err
		}
		err = s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(tx, iconBusinessType, []uint{deviceType.ID})
		if err != nil {
			return err
		}

		// c. 保存新上传的文件
		newFileRecord, err := s.getNewAttachmentRecord(file, adminID, deviceType.ID, fileBusinessType)
		if err != nil {
			return err
		}
		newImageRecord, err := s.getNewAttachmentRecord(image, adminID, deviceType.ID, iconBusinessType)
		if err != nil {
			return err
		}

		// d. 往附件表中写入记录
		err = s.attachmentDao.Create(tx, newFileRecord)
		if err != nil {
			return err
		}
		err = s.attachmentDao.Create(tx, newImageRecord)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("导入设备处理附件提交事务失败: " + err.Error())
	}

	return nil
}

// excelSchema 用于存储从 Header 行解析出的列结构信息
type excelSchema struct {
	// 简单筛选条件: 列索引 -> 条件名称
	Filters map[int]string
	// 范围筛选条件: 起始列索引 -> 条件名称
	RangeFilters map[int]string
	// 组件列: 一个组件的所有列名 -> 列索引 的映射
	ComponentSchema []map[string]int
	// 公共参数列: 列索引 -> 参数名称
	Parameters map[int]string
}

func (s *Service) parseFromExcel(file *multipart.FileHeader) ([]*dto.ImportDataDTO, error) {
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件流失败: %w", err)
	}
	defer f.Close()

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		return nil, fmt.Errorf("读取 Excel 文件失败: %w", err)
	}

	sheetName := xlsx.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("Excel 文件中没有找到任何工作表")
	}

	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("获取 '%s' 工作表数据失败: %w", sheetName, err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("工作表至少需要包含一个标题行和一行数据")
	}

	header := rows[0]
	schema, err := s.buildParsingSchema(xlsx, sheetName, header)
	if err != nil {
		return nil, fmt.Errorf("构建 Excel 解析模式失败: %w", err)
	}
	var allSolutions []*dto.ImportDataDTO
	// 从第二行开始遍历数据
	for _, row := range rows[1:] {
		solutionDTO := &dto.ImportDataDTO{
			Details: &dto.ImportDetailsDTO{
				Filters:    make(map[string]interface{}),
				Components: []*dto.ImportComponentDTO{},
				Parameters: make(map[string]interface{}),
			},
		}

		// 1. 解析筛选条件
		for colIdx, filterName := range schema.Filters {
			if colIdx < len(row) && row[colIdx] != "" {
				solutionDTO.Details.Filters[filterName] = row[colIdx]
			}
		}
		for colIdx, filterName := range schema.RangeFilters {
			if colIdx+1 < len(row) && (row[colIdx] != "" || row[colIdx+1] != "") {
				solutionDTO.Details.Filters[filterName] = map[string]string{
					"min": row[colIdx],
					"max": row[colIdx+1],
				}
			}
		}

		// 2. 解析组件
		for _, componentMap := range schema.ComponentSchema {
			// 检查组件的关键字段（例如商品编码）是否有值，有值才认为是一个有效组件
			keyCode := "商品编码" // 或者 "规格型号"
			keyIndex, ok := componentMap[keyCode]
			if !ok || keyIndex >= len(row) || row[keyIndex] == "" {
				continue // 跳过无效的组件列组
			}

			comp := &dto.ImportComponentDTO{}
			// 使用反射或手动 switch/case 填充字段会更健壮
			if idx, ok := componentMap["工序"]; ok && idx < len(row) {
				comp.Name = row[idx]
			}
			if idx, ok := componentMap["商品编码"]; ok && idx < len(row) {
				comp.ProductCode = row[idx]
			}
			if idx, ok := componentMap["规格型号"]; ok && idx < len(row) {
				comp.SpecCode = row[idx]
			}
			// ... 可以扩展更多组件字段

			solutionDTO.Details.Components = append(solutionDTO.Details.Components, comp)
		}

		// 3. 解析公共参数
		for colIdx, paramName := range schema.Parameters {
			if colIdx < len(row) && row[colIdx] != "" {
				solutionDTO.Details.Parameters[paramName] = row[colIdx]
			}
		}

		allSolutions = append(allSolutions, solutionDTO)
	}

	return allSolutions, nil
}

func (s *Service) buildParsingSchema(xlsx *excelize.File, sheetName string, header []string) (*excelSchema, error) {
	schema := &excelSchema{
		Filters:         make(map[int]string),
		RangeFilters:    make(map[int]string),
		ComponentSchema: []map[string]int{},
		Parameters:      make(map[int]string),
	}

	// 定义颜色
	// 纯色，无透明度。Excelize 返回的是 AARRGGBB 格式，所以我们需要包含 FF 透明度前缀。
	const BlueRgb = "FF0000FF"
	const RedRgb = "FFFF0000"
	const GreenRgb = "FF00FF00" // 纯绿色

	// 识别组件的列名
	componentColumnNames := map[string]bool{
		"工序": true, "商品编码": true, "规格型号": true,
	}

	var currentComponent map[string]int

	for colIdx, colName := range header {
		if colName == "" {
			continue // 跳过空标题
		}

		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		styleID, err := xlsx.GetCellStyle(sheetName, cell)
		if err != nil {
			return nil, fmt.Errorf("获取单元格 '%s' 样式失败: %w", cell, err)
		}

		// 【修正】使用公共 API GetStyle() 获取样式信息
		style, err := xlsx.GetStyle(styleID)
		if err != nil {
			return nil, fmt.Errorf("获取样式 ID %d 的详细信息失败: %w", styleID, err)
		}

		// 检查是否有填充色
		if style.Fill.Type == "pattern" && style.Fill.Pattern > 0 {
			// excelize 返回的颜色通常是 AARRGGBB 格式
			// 我们只取后6位 RGB，并转为大写以便比较
			fgColor := ""
			if len(style.Fill.Color) > 0 {
				// style.Fill.Color 是一个 []string，通常我们只关心第一个
				colorStr := style.Fill.Color[0]
				if len(colorStr) == 8 { // AARRGGBB
					fgColor = strings.ToUpper(colorStr)
				} else if len(colorStr) == 6 { // RGB
					fgColor = "FF" + strings.ToUpper(colorStr) // 补充 alpha 通道
				}
			}

			switch fgColor {
			case BlueRgb:
				schema.Filters[colIdx] = colName
				currentComponent = nil // 中断组件序列
			case RedRgb:
				// 假设范围总是成对出现，我们只记录起始列
				if _, exists := schema.RangeFilters[colIdx-1]; !exists {
					baseName := strings.TrimSuffix(colName, "_min")
					baseName = strings.TrimSuffix(baseName, "_max")
					schema.RangeFilters[colIdx] = baseName
				}
				currentComponent = nil
			case GreenRgb:
				if componentColumnNames[colName] {
					if currentComponent == nil || colName == "工序" {
						currentComponent = make(map[string]int)
						schema.ComponentSchema = append(schema.ComponentSchema, currentComponent)
					}
					currentComponent[colName] = colIdx
				} else {
					schema.Parameters[colIdx] = colName
					currentComponent = nil
				}
			default:
				currentComponent = nil
			}
		} else {
			currentComponent = nil // 无填充色，中断组件序列
		}
	}

	return schema, nil
}

func (s *Service) getNewAttachmentRecord(file *multipart.FileHeader, adminID, businessID uint, businessType string) (*model.Attachment, error) {
	storagePath, err := util.SaveUploadedFile(file)
	if err != nil {
		return nil, err
	}
	newRecord := &model.Attachment{
		Filename:      file.Filename,
		StoragePath:   storagePath,
		FileType:      file.Header.Get("Content-Type"),
		FileSize:      uint64(file.Size),
		StorageDriver: "local",
		UploadedByUID: adminID,
		BusinessType:  util.StringToPointer(businessType),
		BusinessID:    businessID,
	}
	return newRecord, nil
}
