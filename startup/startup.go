package startup

import (
	"news-master/auth"
	"news-master/env"
	"news-master/repository"
	"sync"
)

var loadEnvOnce sync.Once

func Init() {
	loadEnvOnce.Do(_init)
}

func _init() {
	env.LoadEnv()
	auth.LoadKeys()
	repository.Migrate()
}
