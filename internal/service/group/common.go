package group

import (
	"fmt"
	"github.com/spf13/viper"
)

func (s *Service) GetIconMap() (map[uint]string, error) {
	tx := s.dao.DB()
	iconMap := make(map[uint]string)

	baseURL := viper.GetString("server.base_url")
	uploadUrlPrefix := viper.GetString("attachment.upload_url_prefix")
	businessType := viper.GetString("business_type.group_icon")
	icons, err := s.attachmentDao.GetAttachmentsByBusinessType(tx, businessType)
	if err != nil {
		return nil, err
	}
	for _, icon := range icons {
		iconMap[icon.BusinessID] = fmt.Sprintf("%s%s/%s", baseURL, uploadUrlPrefix, icon.StoragePath)
	}
	return iconMap, nil
}
