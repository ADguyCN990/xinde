package group

import (
	"github.com/gin-gonic/gin"
	"xinde/internal/service/group"
)

type Controller struct {
	Service *group.Service
}

func NewGroupController() (*Controller, error) {
	service, err := group.NewGroupService()
	if err != nil {
		return nil, err
	}
	return &Controller{
		Service: service,
	}, nil
}

func (ctrl *Controller) Create(c *gin.Context) {

}
