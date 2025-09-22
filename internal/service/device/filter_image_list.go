package device

import (
	"fmt"
	"github.com/spf13/viper"
	dto "xinde/internal/dto/device"
	"xinde/pkg/stderr"
)

func (s *Service) FilterImageList(deviceTypeID uint, page, pageSize int) (*dto.FilterImageListPageData, error) {
	// 1. 获取总数和总页数
	total, err := s.dao.CountFilterImage(s.dao.DB(), deviceTypeID)
	if err != nil {
		return nil, err
	}
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
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

	// 2. 获取分页后的原始数据
	rawList, err := s.dao.GetFilterImagePage(s.dao.DB(), currentPage, pageSize, deviceTypeID)
	if err != nil {
		return nil, err
	}
	if len(rawList) == 0 {
		return &dto.FilterImageListPageData{
			List:     nil,
			Total:    0,
			Page:     currentPage,
			PageSize: pageSize,
			Pages:    pages,
		}, nil
	}

	// 3. 获取deviceType
	deviceType, err := s.dao.GetDeviceTypeByID(s.dao.DB(), deviceTypeID)
	if err != nil {
		return nil, err
	}

	// 4. 批量查询元数据
	var filterImageIDs []uint
	for _, raw := range rawList {
		filterImageIDs = append(filterImageIDs, uint(raw.ID))
	}
	businessType := viper.GetString("business_type.filter_image")
	images, err := s.attachmentDao.GetAttachmentsByBusinessAndIDs(s.attachmentDao.DB(), businessType, filterImageIDs)
	if err != nil {
		return nil, err
	}
	imageURLMap := make(map[uint]string)
	baseURL := viper.GetString("server.base_url")
	urlPrefix := viper.GetString("attachment.upload_url_prefix")
	for _, img := range images {
		imageURLMap[img.BusinessID] = fmt.Sprintf("%s%s/%s", baseURL, urlPrefix, img.StoragePath)
	}

	// 5. 组装数据
	var listData []*dto.FilterImageListData
	for _, raw := range rawList {
		listData = append(listData, &dto.FilterImageListData{
			ID:             raw.ID,
			DeviceTypeName: deviceType.Name,
			FilterValue:    raw.FilterValue,
			ImageURL:       imageURLMap[raw.ID],
		})
	}
	listPageData := &dto.FilterImageListPageData{
		List:     listData,
		Total:    total,
		Page:     currentPage,
		PageSize: pageSize,
		Pages:    pages,
	}

	// 针对用户输入page过大的情况做特殊处理，返回最后一页的数据，但依然提交err
	if page > pages {
		return listPageData, fmt.Errorf(stderr.ErrorOverLargePage)
	}
	// 针对用户输入page过小的情况做特殊处理，返回第一页的数据，但依然提交error
	if page < 1 {
		return listPageData, fmt.Errorf(stderr.ErrorOverSmallPage)
	}

	return listPageData, nil
}
