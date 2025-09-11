package attachment

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	dto "xinde/internal/dto/attachment"
	model "xinde/internal/model/attachment"
	"xinde/pkg/util"
)

func (s *Service) ScanInvalid() (*dto.OrphanData, error) {
	tx := s.dao.DB()
	savePath := viper.GetString("attachment.save_path")
	if savePath == "" {
		return nil, fmt.Errorf("savePath未配置")
	}
	// 查找孤儿记录，数据库有磁盘没有
	var orphanRecords []*dto.OrphanRecord
	allDBAttachments, err := s.dao.GetAllAttachments(tx)
	if err != nil {
		return nil, err
	}
	// 创建一个记录数据库路径的map，用于查找比对
	dbPathMap := make(map[string]bool, len(allDBAttachments))

	for _, attachment := range allDBAttachments {
		dbPathMap[attachment.StoragePath] = true

		// 检查对应的物理文件是否存在
		absolutePath := filepath.Join(savePath, attachment.StoragePath)
		if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
			// 文件不存在，这是一个孤儿记录
			orphanRecords = append(orphanRecords, converAttachmentToDTOOrphanRecords(attachment))
		}
	}

	//查找孤儿文件，数据库没有磁盘有
	var orphanFiles []string
	err = filepath.Walk(savePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过目录本身
		if !info.IsDir() {
			// 获取文件相对于根目录的存储路径
			relativePath, err := filepath.Rel(savePath, path)
			if err != nil {
				return fmt.Errorf("计算文件相对路径失败")
			}
			// 通过之前创建的map，检查这个文件是否在数据库中存在
			if _, exists := dbPathMap[relativePath]; !exists {
				orphanFiles = append(orphanFiles, relativePath)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("遍历附件目录失败")
	}

	// 组装数据
	result := &dto.OrphanData{
		OrphanRecords: orphanRecords,
		OrphanFiles:   orphanFiles,
	}
	return result, nil
}

func converAttachmentToDTOOrphanRecords(attachment *model.Attachment) *dto.OrphanRecord {
	return &dto.OrphanRecord{
		ID:          attachment.ID,
		Filename:    attachment.Filename,
		StoragePath: attachment.StoragePath,
		UploadedBy:  attachment.UploaderName,
		CreatedAt:   util.FormatTimeToStandardString(attachment.CreatedAt),
	}
}
