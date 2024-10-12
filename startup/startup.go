package startup

import (
	"news-master/auth"
	"news-master/env"
	"sync"
)

var loadEnvOnce sync.Once

func Init() {
	loadEnvOnce.Do(load)
}

func load() {
	env.LoadEnv()
	auth.LoadKeys()
}
