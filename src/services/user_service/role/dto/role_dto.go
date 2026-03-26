package dto

import "url_shortener_pro/src/common/helper"

type RoleCreateDto struct {
	RoleName    string            `json:"roleName"`
	Permissions helper.JsonObject `json:"permissions"`
}

type UserRoleCreate struct {
	UserId int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
}
