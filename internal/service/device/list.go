package device

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
	dto "xinde/internal/dto/device"
	_ "xinde/internal/model/device"
	"xinde/internal/model/group"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) List(page, pageSize int) (*dto.ListPageData, error) {
	pgTx := s.dao.DB()
	mysqlTx := s.groupDao.DB()
	currentPage, pages, err := s.getCurrentPageAndPages(pgTx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 1. 获取分页后的原始DeviceType数据和总数
	count, rawList, err := s.dao.GetDeviceTypeListPage(pgTx, currentPage, pageSize)
	if err != nil {
		return nil, err
	}

	// 2. 一次性获取所有的分组信息，然后在内存中构建分组路径
	allGroups, err := s.groupDao.GetAll(mysqlTx)
	if err != nil {
		return nil, err
	}
	groupMap := make(map[uint]*group.Group)
	for _, g := range allGroups {
		groupMap[g.ID] = g
	}
	// 创建一个缓存，避免重复计算路径
	pathCache := make(map[uint]string)

	// 3.胶合层 收集DeviceTypeID列表，用于批量查询图片
	var deviceTypeIDs []uint
	for _, deviceType := range rawList {
		deviceTypeIDs = append(deviceTypeIDs, deviceType.ID)
	}
	// 4. 从MySQL批量查询DeviceType图片
	businessType := viper.GetString("business_type.device_icon")
	images, err := s.attachmentDao.GetAttachmentsByBusinessAndIDs(mysqlTx, businessType, deviceTypeIDs)
	imageUrlMap := make(map[uint]string)
	baseURL := viper.GetString("server.base_url")
	uploadUrlPrefix := viper.GetString("attachment.upload_url_prefix")
	for _, img := range images {
		imageUrlMap[img.BusinessID] = fmt.Sprintf("%s%s/%s", baseURL, uploadUrlPrefix, img.StoragePath)
	}

	// 5. 组装数据
	var listData []*dto.ListData
	for _, item := range rawList {
		listData = append(listData, &dto.ListData{
			ID:            item.ID,
			GroupName:     s.buildGroupPath(item.GroupID, groupMap, pathCache),
			Name:          item.Name,
			ImageURL:      imageUrlMap[item.ID],
			SolutionCount: item.SolutionCount,
			CreatedAt:     util.FormatTimeToStandardString(item.CreatedAt),
			UpdatedAt:     util.FormatTimeToStandardString(item.UpdatedAt),
		})
	}

	// 6. 组装分页数据
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

func (s *Service) getCurrentPageAndPages(tx *gorm.DB, page, pageSize int) (int, int, error) {
	// 计算页数
	deviceTypeCounts, err := s.dao.CountDeviceTypes(tx)
	if err != nil {
		return 0, 0, err
	}
	pages := int((deviceTypeCounts + int64(pageSize-1)) / int64(pageSize))
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
	return currentPage, pages, nil
}

// buildGroupPath 是一个带缓存的辅助函数，用于在内存中构建分组的完整层级路径
func (s *Service) buildGroupPath(groupID uint, groupMap map[uint]*group.Group, pathCache map[uint]string) string {
	// 如果缓存中已有，直接返回
	if path, ok := pathCache[groupID]; ok {
		return path
	}

	// 使用迭代（循环）代替递归，更安全高效
	var pathParts []string
	currentID := groupID

	for {
		g, ok := groupMap[currentID]
		if !ok {
			// 如果在 map 中找不到，说明数据有问题，中断循环
			break
		}

		// 规则：root 分组 (ID=1) 的名称不加入路径
		if g.ID != 1 {
			pathParts = append(pathParts, g.Name)
		}

		// 如果到达 root (parent_id=0) 或顶级分组 (parent_id=1)，则停止回溯
		if g.ParentID == 0 || g.ParentID == 1 {
			break
		}

		currentID = g.ParentID
	}

	// 反转路径片段
	// pathParts 现在是 [GroupName, level2, level1]，需要反转
	for i, j := 0, len(pathParts)-1; i < j; i, j = i+1, j-1 {
		pathParts[i], pathParts[j] = pathParts[j], pathParts[i]
	}

	fullPath := strings.Join(pathParts, "-")

	// 存入缓存
	pathCache[groupID] = fullPath
	return fullPath
}
