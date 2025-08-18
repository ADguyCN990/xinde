package util

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func SaveUploadedFile(fileHeader *multipart.FileHeader) (string, error) {
	// 从配置中获取存储根目录
	savePath := viper.GetString("attachment.save_path")
	if savePath == "" {
		return "", fmt.Errorf("save_path 未配置")
	}

	// 生成一个唯一的文件名防止冲突
	today := time.Now().Format("20060102")
	ext := filepath.Ext(fileHeader.Filename)
	uniqueFileName := uuid.New().String() + ext

	// 构建完整的目录路径和文件路径
	relativePath := filepath.Join(today, uniqueFileName)
	absolutePath := filepath.Join(savePath, relativePath)

	// 创建目标目录
	if err := os.MkdirAll(filepath.Dir(absolutePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 打开原文件
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件流失败: %w", err)
	}
	defer src.Close()

	// 6. 创建目标文件
	dst, err := os.Create(absolutePath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 7. 将源文件内容拷贝到目标文件
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 8. 返回存入数据库的相对路径
	return relativePath, nil
}
