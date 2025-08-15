package company

import (
	"fmt"
	"xinde/internal/dao/company"
	dto "xinde/internal/dto/company"
	model "xinde/internal/model/account"
	"xinde/pkg/jwt"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

type Service struct {
	dao *company.Dao
	jwt *jwt.JWTService
}

func NewCompanyService() (*Service, error) {
	dao, err := company.NewCompanyDao()
	if err != nil {
		return nil, fmt.Errorf("创建 DAO 实例失败: %w", err)
	}

	jwtService := jwt.NewJWTService()

	return &Service{
		dao: dao,
		jwt: jwtService,
	}, nil
}

func (s *Service) GetCompanyList(page, pageSize int) (*dto.ListPageData, error) {
	tx := s.dao.DB()

	// 计算总页数
	count, err := s.dao.CountCompanies(tx)
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

	// 查询数据库获取当前页面的公司数据
	companies, err := s.dao.FindCompanyListWithPagination(tx, currentPage, pageSize)
	if err != nil {
		return nil, err
	}

	// 将model.Company转换成dto.ListData
	var listData []*dto.ListData
	for _, c := range companies {
		listData = append(listData, convertCompanyToDTOListData(c))
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

func convertCompanyToDTOListData(company *model.Company) *dto.ListData {
	mp := map[string]string{
		"price_1": "价格等级1",
		"price_2": "价格等级2",
		"price_3": "价格等级3",
		"price_4": "价格等级4",
	}

	return &dto.ListData{
		ID:         company.ID,
		Name:       company.Name,
		Address:    util.DerefString(company.Address),
		PriceLevel: mp[company.PriceLevel],
		CreatedAt:  util.FormatTimeToStandardString(company.CreatedAt),
	}
}
