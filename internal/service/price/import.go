package price

import (
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
	_ "time"
	attachmentModel "xinde/internal/model/attachment"
	model "xinde/internal/model/price"
	"xinde/pkg/logger"
	"xinde/pkg/util"
	_ "xinde/pkg/util"
)

func (s *Service) ImportPricesFromFile(c *gin.Context, fileHeader *multipart.FileHeader, adminID uint) error {
	// --- 1. 文件存储 ---
	// a. 打开上传的文件流
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer file.Close()

	// b. 将文件保存到临时位置或直接上传到云存储
	storagePath, err := util.SaveUploadedFile(fileHeader)
	if err != nil {
		return fmt.Errorf("保存上传文件失败: %w", err)
	}

	// 重置文件读取指针，以便后续解析
	file.Seek(0, 0)

	// 在t_attachment表中记录这次上传
	attachment := &attachmentModel.Attachment{
		Filename:      fileHeader.Filename,
		StoragePath:   storagePath,
		FileType:      fileHeader.Header.Get("Content-Type"),
		FileSize:      uint64(fileHeader.Size),
		StorageDriver: "local",
		UploadedByUID: adminID,
		BusinessType:  util.StringToPointer("price_import"),
	}
	if err := s.attachmentDao.Create(s.attachmentDao.DB(), attachment); err != nil {
		// 记录日志，但通常不因为这个失败而中断主流程
		logger.Error("记录上传附件信息到数据库失败: " + err.Error())
	}

	// --- 2. 解析Excel ---
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("读取Excel文件失败: %w", err)
	}
	// 获取工作表中的所有行
	sheetList := xlsx.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("excel文件中没有任何工作表")
	}
	firstSheetName := sheetList[0]
	rows, err := xlsx.GetRows(firstSheetName) // 假设数据在 "Sheet1"
	if err != nil {
		return fmt.Errorf("获取 Sheet1 数据失败: %w", err)
	}
	if len(rows) <= 1 {
		return fmt.Errorf("excel 文件为空或只有表头")
	}

	// 在事务中处理数据并入库
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		for i, row := range rows[1:] {
			// 解析每一行数据，进行类型转化和校验
			priceData, err := s.parsePriceRow(row)
			if err != nil {
				return fmt.Errorf("解析第 %d 行数据失败: %w", i+2, err)
			}
			// 调用dao执行更新或插入操作
			err = s.dao.UpsertPrices(tx, priceData)
			if err != nil {
				return fmt.Errorf("导入第 %d 行数据失败: %w", i+2, err)
			}
		}
		return nil
	})
}

func (s *Service) parsePriceRow(row []string) (*model.Price, error) {
	price1, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return nil, fmt.Errorf("价格数字有误，不是数字类型: %w", err)
	}
	price2, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return nil, fmt.Errorf("价格数字有误，不是数字类型: %w", err)
	}
	price3, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return nil, fmt.Errorf("价格数字有误，不是数字类型: %w", err)
	}
	price4, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, fmt.Errorf("价格数字有误，不是数字类型: %w", err)
	}
	return &model.Price{
		ProductCode: row[0],
		Price1:      price1,
		Price2:      price2,
		Price3:      price3,
		Price4:      price4,
		Unit:        row[5],
		SpecCode:    row[6],
	}, nil
}
