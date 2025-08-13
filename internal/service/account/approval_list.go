package account

import (
	"fmt"
	dto "xinde/internal/dto/account"
	model "xinde/internal/model/account"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) GetApprovalUserList(page, pageSize, status int) (*dto.ApprovalListPageData, error) {
	tx := s.dao.DB()

	// 计算总页数
	count, err := s.dao.CountUserWithStatus(tx, status)
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

	// 查询数据库获取当前页面的用户数据
	dbData, err := s.dao.FindUserListWithPagination(tx, currentPage, pageSize, status)
	if err != nil {
		return nil, err
	}

	// 将model.User转换成dto.ApprovalListData
	var listData []*dto.ApprovalListData
	for _, user := range dbData {
		listData = append(listData, convertUserToDTOApprovalListData(user))
	}

	// 组装分页
	pageData := &dto.ApprovalListPageData{
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

func convertUserToDTOApprovalListData(user *model.User) *dto.ApprovalListData {
	return &dto.ApprovalListData{
		ID:          user.UID,
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		Email:       util.DerefString(user.UserEmail),
		CompanyName: user.CompanyName,
		CreatedAt:   util.FormatTimeToStandardString(user.CreatedAt),
		Why:         util.DerefString(user.Why),
	}
}
