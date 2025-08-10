package cors

import (
	"poketier/pkg/str"

	"github.com/gin-contrib/cors"
)

var (
	allowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	allowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
)

func GetCORSConfig(allowOrigins string, appEnv string) cors.Config {
	if appEnv == "production" {
		return getCORSConfigForProd(allowOrigins)
	} else {
		return getCORSConfigForLocal()
	}
}

// getCORSConfigForProd はプロダクション環境のCORS設定を取得します。
func getCORSConfigForProd(allowOrigins string) cors.Config {
	return cors.Config{
		AllowOrigins:     str.CommaSeparatedToSlice(allowOrigins),
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
	}
}

// getCORSConfigForLocal はローカル環境のCORS設定を取得します。
func getCORSConfigForLocal() cors.Config {
	return cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
	}
}
