package price

import (
	"github.com/gin-gonic/gin"
	"xinde/internal/service/price"
)

type Controller struct {
	priceService *price.Service
}

func NewController() (*Controller, error) {
	priceService, err := price.NewPriceService()
	if err != nil {
		return nil, err
	}

	return &Controller{
		priceService: priceService,
	}, nil
}

func (ctrl *Controller) List(c *gin.Context) {
	
}
