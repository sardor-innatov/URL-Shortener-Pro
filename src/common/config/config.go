package config

import "fmt"

func (c Config) GetDbAdress() string {

	port := c.port
	host := c.host
	res := fmt.Sprintf("%s:%d", host, port)
	return res
}