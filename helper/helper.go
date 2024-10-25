package helper

import (
	"fmt"
	"news-master/app"
)

func PreAuthLink(token string) string {
	return fmt.Sprintf("%v/#/preferences?authToken=%v", app.Config.SiteUrl, token)
}
