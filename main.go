package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/airdb/sailor/dbutil"
	"github.com/airdb/sailor/deployutil"
	"github.com/airdb/sailor/faas"
	"github.com/airdb/sls-bbhj/internal/app"
	"github.com/airdb/sls-bbhj/internal/controller"
	"github.com/airdb/sls-bbhj/internal/repository/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // docs is generated by Swag CLI, you have to import it.
)

const project = "mina"

// @title Airdb Serverlesss Example API
// @version 0.0.1
// @description 宝贝回家小程序后端, 遵循 restful api 规范.
// @termsOfService https://airdb.github.io/terms.html

// @contact.name APIs Support
// @contact.url https://work.weixin.qq.com/kfid/kfc02343d9ba414880a?sence=swagger
// @contact.email info@airdb.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @info.x-logo.url: http://www.apache.org/licenses/LICENSE-2.0.html

// @host service-iw6drlfr-1251018873.sh.apigw.tencentcs.com
// @BasePath /test/mina
// @schemes https

func main() {
	log.Println("start serverless:", deployutil.GetDeployStage(), os.Environ())

	app.InitApp()

	mysqlRepo, err := mysql.GetFactoryOr(dbutil.WriteDefaultDB())
	if err != nil {
		log.Panic(err)
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	// p := filepath.Join("/", deployutil.GetDeployStage(), "/", project)
	p := filepath.Join("/", project)

	mux.Route(p, func(r chi.Router) {
		r.Get("/version", faas.HandleVersion)

		r.Get("/wechat/check_session", controller.CheckSession)

		categoryController := controller.NewCategoryController(mysqlRepo)
		r.Mount("/v1/category", categoryController.Routes())

		lostController := controller.NewLostController(mysqlRepo)
		r.Mount("/v1/lost", lostController.Routes())

		rescueController := controller.NewRescueController(mysqlRepo)
		r.Mount("/v1/rescue", rescueController.Routes())
	})

	if os.Getenv("RUN_MODE") == "local" {
		http.ListenAndServe(":3333", mux)
	} else {
		faas.RunTencentChiWithSwagger(mux, project)
	}
}
