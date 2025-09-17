package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mime/multipart"
	deviceModel "xinde/internal/model/device"
	"xinde/pkg/stderr"
)

func (s *Service) UpdateImport(deviceTypeID, adminID uint, file *multipart.FileHeader) error {

	// 1. 解析Excel。这一步只做纯粹的解析，不涉及任何数据库或API调用。
	parsedData, err := s.parseFromExcel(file)
	if err != nil {
		return err
	}
	if len(parsedData) == 0 {
		return fmt.Errorf("excel没有解析到有效内容")
	}

	// 2. 开启Postgres事务
	err = s.dao.DB().Transaction(func(tx *gorm.DB) error {

		// a. 确认要更新的DeviceType是否存在
		_, err := s.dao.GetDeviceTypeByID(tx, deviceTypeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(stderr.ErrorDeviceNotFound)
			} else {
				return err
			}
		}

		// b. 删除与此DeviceTypeID关联的所有旧方案
		err = s.dao.DeleteByDeviceTypeID(tx, deviceTypeID)
		if err != nil {
			return err
		}

		// c. 准备批量创建的新方案数据
		var solutions []*deviceModel.Device
		for i, data := range parsedData {
			deviceName := fmt.Sprintf("方案%d", i+1)
			detailJson, err := json.Marshal(data.Details)
			if err != nil {
				return fmt.Errorf("序列化方案的Detail失败: " + err.Error())
			}
			solution := &deviceModel.Device{
				Name:         deviceName,
				DeviceTypeID: deviceTypeID,
				Details:      detailJson,
			}
			solutions = append(solutions, solution)
		}

		// d. 批量写入新的 "方案" (Device) 到数据库
		if len(solutions) > 0 {
			err := s.dao.BatchCreateDevice(tx, solutions)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if err.Error() == stderr.ErrorDeviceNotFound {
			return err
		}
		return fmt.Errorf("导入设备提交事务失败: " + err.Error())
	}

	// 3. 处理上传的Excel附件
	err = s.attachmentDao.DB().Transaction(func(tx *gorm.DB) error {

		// a. 获取businessType
		businessType := viper.GetString("business_type.device_import")

		// b. 查找并删除旧的数据库记录
		err := s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(tx, businessType, []uint{deviceTypeID})
		if err != nil {
			return err
		}

		// c. 保存新上传的文件
		newFileRecord, err := s.getNewAttachmentRecord(file, adminID, deviceTypeID, businessType)
		if err != nil {
			return err
		}

		// d. 往附件表中写入记录
		err = s.attachmentDao.Create(tx, newFileRecord)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("导入设备提交附件处理事务失败: " + err.Error())
	}

	return nil
}
