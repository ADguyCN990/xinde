package device

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	dto "xinde/internal/dto/device"
	"xinde/pkg/stderr"
)

func (s *Service) GroupDeviceList(groupID uint) ([]*dto.GroupDeviceListData, error) {
	// 1. 判断是否存在该Group
	group, err := s.groupDao.GetGroupByID(s.groupDao.DB(), groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(stderr.ErrorGroupNotFound)
		} else {
			return nil, err
		}
	}

	// 2. 根据GroupID获取DeviceType列表
	deviceTypeList, err := s.dao.GetDeviceTypesByGroupID(s.dao.DB(), groupID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 3. 根据DeviceTypeID往attachment表里获取图片
	imageMap := make(map[uint]string)
	var deviceIDList []uint

	for _, deviceType := range deviceTypeList {
		deviceIDList = append(deviceIDList, deviceType.ID)
	}
	businessType := viper.GetString("business_type.device_icon")
	images, err := s.attachmentDao.GetAttachmentsByBusinessAndIDs(s.attachmentDao.DB(), businessType, deviceIDList)
	if err != nil {
		return nil, err
	}

	baseURL := viper.GetString("server.base_url")
	uploadUrlPrefix := viper.GetString("attachment.upload_url_prefix")
	for _, image := range images {
		imageMap[image.BusinessID] = fmt.Sprintf("%s%s/%s", baseURL, uploadUrlPrefix, image.StoragePath)
	}

	// 4. 拼接数据
	var list []*dto.GroupDeviceListData
	for _, deviceType := range deviceTypeList {
		list = append(list, &dto.GroupDeviceListData{
			ID:         deviceType.ID,
			DeviceName: deviceType.Name,
			ImageURL:   imageMap[deviceType.ID],
			GroupName:  group.Name,
		})
	}

	return list, nil
}
