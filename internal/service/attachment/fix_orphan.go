package attachment

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	model "xinde/internal/model/attachment"
	"xinde/pkg/logger"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) FixOrphan(adminID uint, filePath string, action string) error {
	tx := s.dao.DB()

	savePath := viper.GetString("attachment.save_path")
	if savePath == "" {
		return fmt.Errorf("savePath未配置")
	}
	absolutePath := filepath.Join(savePath, filePath)

	// 首先检查下文件是否真的存在于磁盘上
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(stderr.ErrorAttachmentNotFoundOnDesk)
		}
		return fmt.Errorf("访问附件状态失败: " + err.Error())
	}

	switch action {
	case "sync":
		// 在数据库中追加一条对应的记录
		fileType := "异常文件，暂不支持该字段"
		businessType := util.StringToPointer("synced_from_orphan")
		attachment := &model.Attachment{
			Filename:      fileInfo.Name(),
			StoragePath:   filePath,
			FileType:      fileType,
			FileSize:      uint64(fileInfo.Size()),
			StorageDriver: "local",
			UploadedByUID: adminID,
			BusinessType:  businessType,
		}

		if err := s.dao.Create(tx, attachment); err != nil {
			// 记录日志，但通常不因为这个失败而中断主流程
			logger.Error("记录上传附件信息到数据库失败: " + err.Error())
		}
	case "delete":
		// 从磁盘上删除文件
		if err := os.Remove(absolutePath); err != nil {
			return fmt.Errorf("删除孤儿文件失败: " + err.Error())
		}
	}
	return nil
}
