package price

import (
	"fmt"
	"xinde/internal/dao/price"
	dto "xinde/internal/dto/price"
	model "xinde/internal/model/price"
	"xinde/pkg/jwt"
	"xinde/pkg/stderr"
)

type Service struct {
	dao *price.Dao
	jwt *jwt.JWTService
}

func NewPriceService() (*Service, error) {
	dao, err := price.NewPriceDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: %v", err)
	}
	jwtService := jwt.NewJWTService()
	return &Service{
		dao: dao,
		jwt: jwtService,
	}, nil
}

func (s *Service) GetPriceList(page, pageSize int) (*dto.ListPageData, error) {
	tx := s.dao.DB()

	// 计算页数
	count, err := s.dao.CountPrices(tx)
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

	// 查询数据库获取当前页面的价格数据
	priceList, err := s.dao.FindPriceListWithPagination(tx, currentPage, pageSize)
	if err != nil {
		return nil, err
	}

	// 将model.Price转换成dto.ListData
	var list []*dto.ListData
	for _, p := range priceList {
		list = append(list, convertPriceToDTOListData(p))
	}

	// 组装分页数据
	pageData := &dto.ListPageData{
		List:     list,
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

func convertPriceToDTOListData(price *model.Price) *dto.ListData {
	return &dto.ListData{
		ID:          price.ID,
		ProductCode: price.ProductCode,
		Price1:      price.Price1,
		Price2:      price.Price2,
		Price3:      price.Price3,
		Price4:      price.Price4,
		Unit:        price.Unit,
		SpecCode:    price.SpecCode,
	}
}
