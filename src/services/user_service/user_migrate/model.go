package user_migrate

import (
	userModel "url_shortener_pro/src/services/user_service/user/model"
	roleModel "url_shortener_pro/src/services/user_service/role/model"
)

func Models() []any {
	return []any{
		userModel.User{},
		roleModel.Role{},
		roleModel.Permission{},
		roleModel.UserRole{},
	}
}