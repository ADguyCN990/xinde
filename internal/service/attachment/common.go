package attachment

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	model "xinde/internal/model/attachment"
)

func (s *Service) getAbsolutePath(attachment *model.Attachment) (string, error) {
	savePath := viper.GetString("attachment.save_path")
	if savePath == "" {
		return "", fmt.Errorf("savePath未配置")
	}
	absolutePath := filepath.Join(savePath, attachment.StoragePath)
	return absolutePath, nil
}
