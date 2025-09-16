package device

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
	"xinde/internal/dao/attachment"
	"xinde/internal/dao/device"
	model "xinde/internal/model/attachment"
	"xinde/pkg/jwt"
	"xinde/pkg/util"
)

type Service struct {
	dao           *device.Dao
	j             *jwt.JWTService
	attachmentDao *attachment.Dao
}

func NewDeviceService() (*Service, error) {
	dao, err := device.NewDeviceDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	attachmentDao, err := attachment.NewAttachmentDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	j := jwt.NewJWTService()
	return &Service{
		dao:           dao,
		j:             j,
		attachmentDao: attachmentDao,
	}, nil
}

func (s *Service) ImportFromExcel(adminID, groupID uint, deviceTypeName string, file, image *multipart.FileHeader) error {
	tx := s.attachmentDao.DB()

	// 将导入的文件记录到附件管理中。暂不create因为缺少businessID
	newImageRecord, err := getNewAttachmentRecord(image, adminID)
	if err != nil {
		return err
	}
	newExcelRecord, err := getNewAttachmentRecord(file, adminID)
	if err != nil {
		return err
	}
}

func getNewAttachmentRecord(file *multipart.FileHeader, adminID uint) (*model.Attachment, error) {
	// 上传文件到磁盘，暂不写入attachment表（还不知道device_type的ID）
	storagePath, err := util.SaveUploadedFile(file)
	if err != nil {
		return nil, err
	}
	businessType := viper.GetString("business_type.device_import")
	newRecord := &model.Attachment{
		Filename:      file.Filename,
		StoragePath:   storagePath,
		FileType:      file.Header.Get("Content-Type"),
		FileSize:      uint64(file.Size),
		StorageDriver: "local",
		UploadedByUID: adminID,
		BusinessType:  util.StringToPointer(businessType),
		BusinessID:    0,
	}
	return newRecord, nil
}
