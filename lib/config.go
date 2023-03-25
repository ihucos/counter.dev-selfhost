package lib

import (
	"fmt"
	rand "crypto/rand"
	"encoding/hex"
	"os"
	"time"
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	RedisUrl     string
	Bind         string
	CookieSecret []byte
	PasswordSalt []byte
}

func env(env string) string {
	v := os.Getenv(env)
	if v == "" {
		panic(fmt.Sprintf("empty or missing env: %s", env))
	}
	return v
}
func envDefault(env string, fallback string) string {
	v := os.Getenv(env)
	if v == "" {
		return fallback
	}
	return v
}


func envDuration(envName string) time.Duration {
	strVal := env(envName)
	duration, err := time.ParseDuration(strVal)
	if err != nil {
		panic(fmt.Sprintf("Not duration given for: %s; %s", envName, err))
	}
	return duration
}

func getSecret(conn redis.Conn, key string) string{

	// generate secret token
	tokenLength := 64 
	tokenBytes := make([]byte, tokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic(err)
	}
	newSecret := hex.EncodeToString(tokenBytes)

	// set new secret only if it does not exist yet
	_, err = redis.String(conn.Do("SET", key, newSecret, "NX"))
	if err != redis.ErrNil {
		panic(err)
	}

	// get the saved secret token
	secret, err := redis.String(conn.Do("GET", key))
	if err != nil {
		panic(err)
	}
	return secret
}

func NewConfig() Config {
	redisUrl := envDefault("COUNTER_REDIS_URL", "redis://localhost:6379")
	conn, err := redis.DialURL(redisUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	return Config{
		RedisUrl:     redisUrl,
		Bind:         envDefault("COUNTER_BIND", ":8080"),
		CookieSecret: []byte(getSecret(conn, "cntr:config:cookie_secret")),
		PasswordSalt: []byte(getSecret(conn, "cntr:config:password_salt")),
	}

}
