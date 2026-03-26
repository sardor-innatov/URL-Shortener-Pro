package seeder

import (
	role_model "url_shortener_pro/src/services/user_service/role/model"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {

	permissions := []role_model.Permission{

		// auth endpoints
		{Path: "/api/v1/auth", EndpointName: "sign up", Method: "POST"},
		{Path: "/api/v1/auth/login", EndpointName: "login", Method: "POST"},
		{Path: "/api/v1/auth/me", EndpointName: "get users info", Method: "GET"},
		{Path: "/api/v1/auth/refresh", EndpointName: "refresh token", Method: "POST"},
		// role endpoints
		{Path: "/api/v1/role", EndpointName: "get all roles", Method: "GET"},
		{Path: "/api/v1/role/permissions", EndpointName: "get all permissions", Method: "GET"},
		{Path: "/api/v1/role", EndpointName: "create role", Method: "POST"},
		{Path: "/api/v1/role/:id", EndpointName: "update role", Method: "PUT"},
		{Path: "/api/v1/role/:id", EndpointName: "get role by id", Method: "GET"},
		{Path: "/api/v1/role/:id", EndpointName: "delete role", Method: "DELETE"},
		// link endpoints
		{Path: "/api/v1/link/delete/:id", EndpointName: "delete link", Method: "DELETE"},
		{Path: "/api/v1/link/my", EndpointName: "get user links", Method: "GET"},
		{Path: "/api/v1/link/shorten", EndpointName: "create short link", Method: "POST"},
		{Path: "/:shortCode", EndpointName: "redirect", Method: "GET"},
		// stats endpoints
		{Path: "/api/v1/link/:id/stats", EndpointName: "get stats", Method: "GET"},
	}

	for _, p := range permissions {
		// using Where to not dublicate permissions every boot
		// Используем Where, чтобы не дублировать права при каждом запуске
		db.Where(role_model.Permission{Path: p.Path, Method: p.Method}).
			FirstOrCreate(&p)
	}


	
}
