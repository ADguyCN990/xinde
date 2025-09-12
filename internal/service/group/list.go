package group

import (
	"fmt"
	dto "xinde/internal/dto/group"
	model "xinde/internal/model/group"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) GetGroupList(page, pageSize int) (*dto.ListPageData, error) {
	// 获取所有分组，为了计算层级这是必须的
	tx := s.dao.DB()
	allGroups, err := s.dao.GetAll(tx)
	if err != nil {
		return nil, err
	}

	// 在内存中计算每个节点的层级
	groupMap := make(map[uint]*model.Group, len(allGroups))
	levelMap := make(map[uint]int, len(allGroups))
	for _, g := range allGroups {
		groupMap[g.ID] = g
	}
	for _, g := range allGroups {
		levelMap[g.ID] = s.getLevel(g.ID, groupMap, levelMap)
	}

	// 获取图标URL映射
	iconMap, err := s.GetIconMap()
	if err != nil {
		return nil, err
	}

	// 手动计算分页数据
	count := len(allGroups)
	pages := (count + pageSize - 1) / pageSize
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

	// 通过切片手动截取分页数据
	start := (currentPage - 1) * pageSize
	end := start + pageSize
	if end > count {
		end = count
	}
	paginatedData := allGroups[start:end]

	var backendList []*dto.BackendListData
	for _, group := range paginatedData {
		backendList = append(backendList, convertGroupToDTOListPageData(group, groupMap, levelMap, iconMap))
	}

	// 组装分页数据
	pageData := &dto.ListPageData{
		List:     backendList,
		Total:    count,
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

func convertGroupToDTOListPageData(group *model.Group, groupMap map[uint]*model.Group, levelMap map[uint]int, iconMap map[uint]string) *dto.BackendListData {
	parentName := ""
	if group.ParentID != 0 {
		if parent, ok := groupMap[group.ParentID]; ok {
			parentName = parent.Name
		}
	}
	return &dto.BackendListData{
		ID:         group.ID,
		Name:       group.Name,
		ParentID:   group.ParentID,
		ParentName: parentName,
		Level:      levelMap[group.ID],
		IconURL:    iconMap[group.ID],
		CreatedAt:  util.FormatTimeToStandardString(group.CreatedAt),
	}
}

func (s *Service) getLevel(id uint, groupMap map[uint]*model.Group, levelMap map[uint]int) int {
	if level, ok := levelMap[id]; ok {
		return level
	}
	g := groupMap[id]
	if g.ParentID == 0 {
		// 找到了root节点，root节点的层级为0
		levelMap[g.ID] = 0
		return 0
	}

	// 对于子节点，向上递归父节点的level，自己的level就是父节点的level+1
	parentLevel := s.getLevel(g.ParentID, groupMap, levelMap)
	level := parentLevel + 1
	levelMap[g.ID] = level
	return level
}
