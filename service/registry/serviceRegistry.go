package registry

import (
	"strings"
	"time"
)

func Register(serviceName, ip, host string, timeLen time.Duration) {
	go func() {
		for {
			//TODO 保留getLocalUrl接口，之后会改成通过该函数获得当前的ip地址
			url := getLocalUrl(ip, host)
			rds.SAdd(ctx, serviceName, url)
			rds.Set(ctx, url, "ok", timeLen+time.Second)
			time.Sleep(timeLen)
		}
	}()
}

func getLocalUrl(ip, host string) string {
	var LocalUrl strings.Builder
	LocalUrl.WriteString(ip)
	LocalUrl.WriteString(":")
	LocalUrl.WriteString(host)
	return LocalUrl.String()
}
