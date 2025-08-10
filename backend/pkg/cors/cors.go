package cors

import (
	"poketier/pkg/str"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	allowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	allowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
)

func SetCORS(r *gin.Engine, allowOrigins string, appEnv string) {
	if appEnv == "production" {
		setCORSForProd(r, allowOrigins)
	} else {
		setCORSForLocal(r)
	}
}

// setCORSForProd は CORS 設定を適用します。プロダクション環境では環境変数で設定された特定のオリジンのみを許可します。
func setCORSForProd(r *gin.Engine, allowOrigins string) {
	corsConfig := cors.Config{
		AllowOrigins:     str.CommaSeparatedToSlice(allowOrigins),
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))
}

// setCORSForLocal は CORS 設定を適用します。ローカル環境では全てのオリジンを許可します。
func setCORSForLocal(r *gin.Engine) {
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))
}
