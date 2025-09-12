package group

import (
	dto "xinde/internal/dto/group"
)

func (s *Service) GetTree(includedIcon string) ([]*dto.GroupTreeNode, error) {
	tx := s.dao.DB()
	allGroups, err := s.dao.GetAll(tx)
	if err != nil {
		return nil, err
	}

	// 如果需要图标，还要获取所有分组对应的图标映射
	iconMap := make(map[uint]string)
	if includedIcon == "true" {
		iconMap, err = s.GetIconMap()
		if err != nil {
			return nil, err
		}
	}

	// 将扁平化列表转换为树状结构
	// 第一次遍历先用map做预处理
	nodeMap := make(map[uint]*dto.GroupTreeNode)
	for _, group := range allGroups {
		node := &dto.GroupTreeNode{
			ID:       group.ID,
			Name:     group.Name,
			ParentID: group.ParentID,
			Children: []*dto.GroupTreeNode{}, //初始化为空切片，避免json序列化为null
		}
		if includedIcon == "true" {
			if url, ok := iconMap[group.ID]; ok {
				node.IconURL = url
			}
		}
		nodeMap[group.ID] = node
	}

	var tree []*dto.GroupTreeNode
	for _, node := range nodeMap {
		if node.ParentID == 0 {
			// 根节点root
			tree = append(tree, node)
		} else {
			if parent, ok := nodeMap[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}
	return tree, nil
}
