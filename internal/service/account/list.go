package account

import (
	"math"
	_ "xinde/internal/dao/account"
	dto "xinde/internal/dto/account"
	model "xinde/internal/model/account"
	_ "xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) GetUserList(page, pageSize int) (*dto.ListPageData, error) {
	tx := s.dao.DB()

	// 查询数据库获取当前页面的用户数据
	dbData, totalCount, err := s.dao.FindUserListWithPagination(tx, page, pageSize)

	// 没有需要特别注意的业务error，直接返回就行
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
		Total:    totalCount,
		Page:     page,
		PageSize: pageSize,
		Pages:    int(math.Ceil(float64(totalCount) / float64(pageSize))),
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
	return &dto.ListData{
		ID:             user.UID,
		Name:           user.Name,
		Username:       user.Username,
		Phone:          user.Phone,
		Email:          util.DerefString(user.UserEmail),
		CompanyName:    user.CompanyName,
		PriceLevel:     "TODO,价格管理",
		Remark:         util.DerefString(user.Remarks),
		Role:           userRole,
		CreatedAt:      util.FormatNullableTimeToStandardString(user.HandledAt),
		RecentSearchAt: util.FormatNullableTimeToStandardString(user.RecentSearchAt),
		SearchDevice:   util.DerefString(user.SearchDevice),
	}
}
