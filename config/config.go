package config

import (
	"encoding/json"
	"errors"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/george012/gtbox/gtbox_redis"
	"github.com/wmyeah/yeah_box/api/api_config"
	"os"
	"path/filepath"
)

const (
	ProjectName             = "yeah_box"
	ProjectVersion          = "v0.0.2"
	ProjectBundleID         = "com.yeah_box.yeah_box"
	netListenAPIPortDefault = 17173
)

type AccountsConfig struct {
	Enabled             bool   `yaml:"enabled" json:"enabled"`
	SyncTimeCycleSecond int64  `yaml:"sync_time_cycle_second" json:"sync_time_cycle_second"`
	Addr                string `yaml:"addr" json:"addr"`
}

type MiningpoolConfig struct {
	SubAccount string `yaml:"sub" json:"sub"`
	Miner      string `yaml:"miner" json:"miner"`
	Pool       string `yaml:"pool" json:"pool"`
}

type FileConfig struct {
	RedisCfg      gtbox_redis.RedisConfig `yaml:"redis_cfg" json:"redis_cfg"`
	Api           api_config.ApiConfig    `yaml:"api" json:"api"`
	BaseUploadDir string                  `yaml:"upload_dir" json:"upload_dir"`
}

var GlobalConfig *FileConfig

func LoadConfig(file string) error {
	fInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fInfo.IsDir() {
		return errors.New("config file can not be a dir")
	}

	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &GlobalConfig)
	//err = yaml.Unmarshal(buf, &GlobalConfig)
	if err != nil {
		return err
	}

	return nil
}

func SaveConfig(file string) error {
	if file == "" {
		file = CurrentApp.AppConfigFilePath
	}
	//config, err := yaml.Marshal(GlobalConfig)
	config, err := json.MarshalIndent(GlobalConfig, "", "    ")

	if err != nil {
		return err
	}

	err = os.WriteFile(file, config, 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateDefaultConfig() *FileConfig {
	return &FileConfig{
		Api: api_config.ApiConfig{
			Enabled: true,
			Port:    netListenAPIPortDefault,
		},
		BaseUploadDir: "./uploads",
	}
}

func SyncConfigFile() {
	gtbox_log.LogInfof("加载配置文件 [%s]", CurrentApp.AppConfigFilePath)
	_, err := os.Stat(CurrentApp.AppConfigFilePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// 获取配置文件的父目录路径
		dir := filepath.Dir(CurrentApp.AppConfigFilePath)

		// 检查父目录是否存在
		if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			// 创建父目录
			if err = os.MkdirAll(dir, 0755); err != nil {
				gtbox_log.LogErrorf("无法创建目录 [%s]: %s", dir, err.Error())
				return
			}
		}

		// 写入默认配置文件内容
		jd, _ := json.MarshalIndent(generateDefaultConfig(), "", "  ")
		err = os.WriteFile(CurrentApp.AppConfigFilePath, jd, 0755)
		if err != nil {
			gtbox_log.LogErrorf("无法写入配置文件 [%s]: %s", CurrentApp.AppConfigFilePath, err.Error())
			return
		}
	} else {
		buf, err := os.ReadFile(CurrentApp.AppConfigFilePath)
		if err != nil {
			gtbox_log.LogErrorf("读取配置文件 [%s] 错误: %s", CurrentApp.AppConfigFilePath, err.Error())

			return
		}
		if len(buf) == 0 {
			gtbox_log.LogErrorf("配置文件重置")
			jd, _ := json.MarshalIndent(generateDefaultConfig(), "", "  ")
			// 写入默认配置文件内容
			err = os.WriteFile(CurrentApp.AppConfigFilePath, jd, 0755)
			if err != nil {
				gtbox_log.LogErrorf("无法写入配置文件 [%s]: %s", CurrentApp.AppConfigFilePath, err.Error())
				return
			}
		}
	}

	err = LoadConfig(CurrentApp.AppConfigFilePath)

	if err != nil {
		gtbox_log.LogErrorf("无法加载配置文件 [%s]: %s", CurrentApp.AppConfigFilePath, err.Error())
		return
	}

}