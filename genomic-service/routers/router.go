// Source code template generated by https://gitlab.com/technixo/tnx

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	golibgin "github.com/golibs-starter/golib-gin"
	"github.com/golibs-starter/golib/web/actuator"
	"github.com/minhngoc274/genomic-system/genomic-service/controllers"
	"go.uber.org/fx"
)

// RegisterRoutersIn represents constructor params for fx
type RegisterRoutersIn struct {
	fx.In
	App            *golib.App
	Engine         *gin.Engine
	Actuator       *actuator.Endpoint
	UserController *controllers.UserController
}

// RegisterHandlers register handlers
func RegisterHandlers(app *golib.App, engine *gin.Engine) {
	engine.Use(golibgin.InitContext())
	engine.Use(golibgin.WrapAll(app.Handlers())...)
}

// RegisterGinRouters register gin routes
func RegisterGinRouters(p RegisterRoutersIn) {
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))

	group.POST("v1/users", p.UserController.Register)
	group.POST("v1/users/upload", p.UserController.UploadData)
}
