package attachment

import (
	"fmt"
	"xinde/internal/dao/attachment"
	dao "xinde/internal/dao/attachment"
	dto "xinde/internal/dto/attachment"
	model "xinde/internal/model/attachment"
	"xinde/pkg/jwt"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

type Service struct {
	jwt *jwt.JWTService
	dao *attachment.Dao
}

func NewAttachmentService() (*Service, error) {
	jwtService := jwt.NewJWTService()
	dao, err := attachment.NewAttachmentDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	return &Service{
		jwt: jwtService,
		dao: dao,
	}, nil
}

func (s *Service) List(page, pageSize int, filename string) (*dto.ListPageData, error) {
	tx := s.dao.DB()
	params := &dao.ListParams{
		PageSize: pageSize,
		Filename: filename,
	}

	// 计算页数
	count, err := s.dao.CountWithParams(tx, params)
	if err != nil {
		return nil, err
	}
	pages := int((count + int64(pageSize-1)) / int64(pageSize))
	if pages == 0 {
		pages = 1
	}

	// 对page过大的情况做判断
	currentPage := page
	if currentPage > pages {
		currentPage = pages
	}
	// 对page过小的情况做判断
	if currentPage < 1 {
		currentPage = 1
	}
	params.Page = currentPage

	// 查询数据库获取当前页面的附件数据
	list, err := s.dao.FindAttachmentListWithPagination(tx, params)
	if err != nil {
		return nil, err
	}

	// 将model.Attachment转换成dto.ListData
	var listData []*dto.ListData
	for _, a := range list {
		listData = append(listData, convertAttachmentToDTOListData(a))
	}

	// 组装分页数据
	pageData := &dto.ListPageData{
		List:     listData,
		Total:    int(count),
		Page:     currentPage,
		PageSize: pageSize,
		Pages:    pages,
	}

	// 针对用户输入page过大的情况做特殊处理，返回最后一页的数据，但依然提交err
	if page > pages {
		return pageData, fmt.Errorf(stderr.ErrorOverLargePage)
	}
	// 针对用户输入page过小的情况做特殊处理，返回第一页的数据，但依然提交error
	if page < 1 {
		return pageData, fmt.Errorf(stderr.ErrorOverSmallPage)
	}

	return pageData, nil
}

func convertAttachmentToDTOListData(attachment *model.Attachment) *dto.ListData {
	return &dto.ListData{
		ID:            attachment.ID,
		Filename:      attachment.Filename,
		FileType:      attachment.FileType,
		FileSize:      util.FormatFileSize(attachment.FileSize),
		StorageDriver: attachment.StorageDriver,
		BusinessType:  util.DerefString(attachment.BusinessType),
		UploadedBy:    attachment.UploaderName,
		CreatedAt:     util.FormatTimeToStandardString(attachment.CreatedAt),
	}
}
