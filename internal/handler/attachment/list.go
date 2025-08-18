package attachment

import (
	"fmt"
	"xinde/internal/service/attachment"
)

type Controller struct {
	attachmentService *attachment.Service
}

func NewAttachmentController() (*Controller, error) {
	service, err := attachment.NewAttachmentService()
	if err != nil {
		return nil, fmt.Errorf("创建Service实例失败: " + err.Error())
	}
	return &Controller{attachmentService: service}, nil
}
