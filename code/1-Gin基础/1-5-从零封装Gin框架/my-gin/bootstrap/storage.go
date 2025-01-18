package bootstrap

import (
	"github.com/jassue/go-storage/local"
	"my-gin/global"
)

func InitializeStorage() {
	_, _ = local.Init(global.App.Config.Storage.Disks.Local)
}
