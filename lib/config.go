package lib

import (
)

type Config struct {
	RedisUrl     string
	Bind         string
	CookieSecret []byte
	PasswordSalt []byte
}
