package settings

import (
	"github.com/spf13/viper"
)


type Setting struct {
	vp *viper.Viper
}
// 初始配置方法，读取指定配置文件
func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}