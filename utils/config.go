package utils

import (
	"github.com/spf13/viper"
)

func LoadConfigInFile(filename string) error {

	// 如果指定了文件，则从指定文件读取，如果没指定则按默认规则加载
	// 默认规则为：当前目录/{env}.{ext}
	// ext 支持的类型有 json yaml 等
	if filename == "" {
		env := GetEnvironment()

		// 1. 读取配置文件 「app」
		if FilesExists([]string{"config/app.yml", "config/app.yaml", "config/app.json"}) {
			v := viper.New()
			v.AddConfigPath("config")
			v.SetConfigName("app")
			if err := v.ReadInConfig(); err != nil {
				return err
			}

			viper.MergeConfigMap(v.AllSettings())
		}

		// 2. 读取配置文件 「env」 ，会覆盖 app 里的配置
		if FilesExists([]string{"config/" + env + ".yml", "config/" + env + ".yaml", "config/" + env + ".json"}) {
			v := viper.New()
			v.AddConfigPath("config")
			v.SetConfigName(env)
			if err := v.ReadInConfig(); err != nil {
				return err
			}

			viper.MergeConfigMap(v.AllSettings())
		}

		// 3. 读取自定义配置 「custom」，会覆盖环境默认配置
		if FilesExists([]string{"config/custom.yml", "config/custom.yaml", "config/custom.json"}) {
			v := viper.New()
			v.AddConfigPath("config")
			v.SetConfigName("custom")
			if err := v.ReadInConfig(); err != nil {
				return err
			}

			viper.MergeConfigMap(v.AllSettings())
		}

	} else {
		viper.SetConfigFile(filename)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	return nil
}
