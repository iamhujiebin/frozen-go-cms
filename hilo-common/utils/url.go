package utils

import (
	"frozen-go-cms/hilo-common/resource/config"
	"strings"
)

const DefaultAvatarMan = "hilo/manager/ea48b62d54a24a709de3c38702c89995.png"
const DefaultAvatarWoman = "hilo/manager/ea48b62d54a24a709de3c38702c89995.png"
const DEFAULT_NICK = "Hilo No.%s"

// 补全url，区分处理oss和aws两种情况
func MakeFullUrl(url string) string {
	if strings.HasPrefix(url, config.GetConfigOss().OSS_CDN) || strings.HasPrefix(url, config.GetConfigAws().AWS_CDN) {
		return url
	} else if strings.HasPrefix(url, "nextvideo/") {
		return config.GetConfigOss().OSS_CDN + url
	} else if strings.HasPrefix(url, config.GetConfigAws().AWS_DIR) {
		return config.GetConfigAws().AWS_CDN + url
	} else {
		return url
	}
}
