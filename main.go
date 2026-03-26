package main

import cmd "url_shortener_pro/src"

// @title           URL Shortener Pro API
// @version         1.0
// @BasePath  /api/v1
// @host      localhost:8080
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Введите токен в формате: Bearer <your_token>
func main() {
	cmd.Exec()
}