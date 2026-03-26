package cmd

import (
	"fmt"

	"url_shortener_pro/src/common/config"
	"url_shortener_pro/src/common/helper"
	"url_shortener_pro/src/services/link_service/link_migrate"
	"url_shortener_pro/src/services/user_service/user_migrate"

	role_model "url_shortener_pro/src/services/user_service/role/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func configureEnv(router *echo.Echo) config.EnvProject {
	return config.ProjectEnv()
}

func migrate(db *gorm.DB) {

	err := user_migrate.Migrate(db, user_migrate.Models())
	{
		if err != nil {
			panic(err.Error())
		}
	}
	err = link_migrate.Migrate(db, link_migrate.Models())
	{
		if err != nil {
			panic(err.Error())
		}
	}
}

func createOrg(db *gorm.DB, r *echo.Echo) {

	routers := make(helper.JsonObject)
	for _, route := range r.Routes() {
		if route.Method == "echo_route_not_found" {
			// fmt.Println("Skipping route:", route.Path, "Method:", route.Method)
			continue
		}
		if existing, ok := routers[route.Path]; ok {
			routers[route.Path] = append(existing.([]string), route.Method)
		} else {
			routers[route.Path] = []string{route.Method}
		}
		//jsonbPermissions[route.Path] = route.Method
	}
	fmt.Println(db.Where("roles.id = 1 OR roles.id = 2").Delete(&role_model.Role{}).Error)

	fmt.Println(routers)

	//jsonbPermissions = routers
	//fmt.Println(jsonbPermissions)

	roles := []role_model.Role{
		{
			Id:          1,
			RoleName:    "superadmin",
			Permissions: routers,
		},
		{
			Id:       2,
			RoleName: "user",
			Permissions: helper.JsonObject{
				// redirect
				"/:shortCode": []string{"GET"},
				// auth
				"/api/v1/auth":         []string{"POST"},
				"/api/v1/auth/login":   []string{"POST"},
				"/api/v1/auth/me":      []string{"GET"},
				"/api/v1/auth/refresh": []string{"POST"},
				// link
				"/api/v1/link/my":      []string{"GET"},
				"/api/v1/link/shorten": []string{"POST"},
			},
		},
	}
	db.Create(&roles)
}
