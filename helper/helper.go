package helper

import (
	"fmt"
	"news-master/app"
)

func PreAuthLink(token string) string {
	return fmt.Sprintf("%v/#/preferences?authToken=%v", app.Config.SiteUrl, token)
}

func AboutLinkLink() string {
	return fmt.Sprintf("%v/#/about", app.Config.SiteUrl)
}

func PrivacyLink() string {
	return fmt.Sprintf("%v/#/privacy", app.Config.SiteUrl)
}
