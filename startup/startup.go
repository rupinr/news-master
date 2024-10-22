package startup

import (
	"news-master/app"
	"news-master/auth"
	"news-master/repository"
	"sync"
)

var loadEnvOnce sync.Once

func Init() {
	loadEnvOnce.Do(_init)
}

func _init() {
	app.Load()
	auth.LoadKeys()
	repository.Migrate()
}
