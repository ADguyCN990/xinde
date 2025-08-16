package configs

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig(filePath string) error {
	// 1. 检查文件路径是否提供
	if filePath == "" {
		return fmt.Errorf("配置文件路径不能为空")
	}

	// 2. 直接设置配置文件的完整路径
	// 我们不再使用 AddConfigPath, SetConfigName, SetConfigType
	// 而是直接告诉 viper 读取哪个文件
	viper.SetConfigFile(filePath)

	// 3. 启用环境变量覆盖 (可选但推荐)
	// 这允许你通过环境变量来覆盖配置文件中的值，例如:
	// export SERVER_PORT=9090
	// viper 会自动寻找名为 SERVER.PORT 的配置项
	viper.AutomaticEnv()
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 可选，用于将 SERVER.PORT 转换为 SERVER_PORT

	// 4. 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果文件不存在，viper 会返回一个特定的错误类型
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("未在该路径下找到配置文件: %s", filePath)
		}

		// 其他类型的错误，比如 YAML 格式错误
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 如果需要，可以在这里设置一些默认值
	// viper.SetDefault("server.port", 8080)

	return nil
}
