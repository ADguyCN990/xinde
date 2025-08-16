package account

import (
	"fmt"
	_ "math"
	_ "xinde/internal/dao/account"
	dto "xinde/internal/dto/account"
	model "xinde/internal/model/account"
	"xinde/pkg/stderr"
	_ "xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) GetUserList(page, pageSize int) (*dto.ListPageData, error) {
	tx := s.dao.DB()

	// 计算总页数
	count, err := s.dao.CountUserWithStatus(tx, model.UserApproved)
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
	dbData, err := s.dao.FindUserListWithPagination(tx, currentPage, pageSize, model.UserApproved)
	if err != nil {
		return nil, err
	}

	// 将model.User转换成dto.ListData
	var listData []*dto.ListData
	for _, user := range dbData {
		listData = append(listData, s.convertUserToDTOListData(user))
	}

	// 组装分页
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

func (s *Service) convertUserToDTOListData(user *model.User) *dto.ListData {
	var userRole string
	if user.IsAdmin == 1 {
		userRole = "管理员"
	} else {
		userRole = "普通用户"
	}

	//TODO 用户访问记录
	return &dto.ListData{
		ID:             user.UID,
		Name:           user.Name,
		Username:       user.Username,
		Phone:          user.Phone,
		Email:          util.DerefString(user.UserEmail),
		CompanyName:    user.CompanyName,
		PriceLevel:     user.PriceLevel,
		Remark:         util.DerefString(user.Remarks),
		Role:           userRole,
		CreatedAt:      util.FormatNullableTimeToStandardString(user.HandledAt),
		RecentSearchAt: util.FormatNullableTimeToStandardString(user.RecentSearchAt),
		SearchDevice:   util.DerefString(user.SearchDevice),
	}
}
