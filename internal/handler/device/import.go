package device

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xinde/internal/service/device"
)

type Controller struct {
	service *device.Service
}

func NewDeviceController() (*Controller, error) {
	service, err := device.NewDeviceService()
	if err != nil {
		return nil, fmt.Errorf("创建service实例失败: " + err.Error())
	}
	return &Controller{service: service}, nil
}

func (ctrl *Controller) Import(c *gin.Context) {

}
