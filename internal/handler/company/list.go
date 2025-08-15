package company

import (
	"github.com/gin-gonic/gin"
	"xinde/internal/service/company"
	service "xinde/internal/service/company"
)

type Controller struct {
	companyService *company.Service
}

func NewCompanyController() (*Controller, error) {
	companyService, err := service.NewCompanyService()
	if err != nil {
		return nil, err
	}

	return &Controller{
		companyService: companyService,
	}, nil
}

func (ctrl *Controller) List(c *gin.Context) {

}
