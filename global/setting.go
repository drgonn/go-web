package global

import (
	"github.com/go-web/pkg/settings"
	"github.com/go-web/pkg/logger"
)

var (
	ServerSetting   *settings.ServerSettingS
	AppSetting      *settings.AppSettingS
	DatabaseSetting *settings.DatabaseSettingS
	Logger          *logger.Logger
)