package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/jassue/go-storage/storage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"my-gin/config"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
	Log         *zap.Logger
	DB          *gorm.DB
	Redis       *redis.Client
}

var App = new(Application)

func (app *Application) Disk(disk ...string) storage.Storage {
	// 若未传参，默认使用配置文件驱动
	diskName := app.Config.Storage.Default
	if len(disk) > 0 {
		diskName = storage.DiskName(disk[0])
	}
	s, err := storage.Disk(diskName)
	if err != nil {
		app.Log.Error("error to get the information of file: " + string(diskName))
	}
	return s
}
