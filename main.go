package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"syscall"
	rand "crypto/rand"
	"encoding/hex"

	_ "github.com/ihucos/counter.dev/endpoints"
	"github.com/gomodule/redigo/redis"
	"github.com/ihucos/counter.dev/lib"

	"github.com/urfave/cli/v2"
)

//go:embed all:static
var staticFS embed.FS



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



func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Serve app",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "bind",
						Value: ":80",
						Usage: "host:port to bind server to",
					},
					&cli.StringFlag{
						Name:  "redis-url",
						Value: "redis://localhost:6379",
						Usage: "Which redis server to connect to",
					},
				},
				Action: func(cCtx *cli.Context) error {

					// We want to open many file descriptiors for the redis pooling under moderate load to work
					var rLimit syscall.Rlimit
					rLimit.Max = 100307
					rLimit.Cur = 100307
					err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
					if err != nil {
						fmt.Println("Error Setting Rlimit ", err)
					}


					conn, err := redis.DialURL(cCtx.String("redis-url"))
					if err != nil {
						fmt.Printf("Can't connect to redis server at %s - %s\n", cCtx.String("redis-url"), err)
						os.Exit(1)
					}

					counterApp := lib.NewApp(
						lib.Config{
							RedisUrl:     cCtx.String("redis-server"),
							Bind:         cCtx.String("bind"),
							CookieSecret: []byte(getSecret(conn, "cntr:config:cookie_secret")),
							PasswordSalt: []byte(getSecret(conn, "cntr:config:password_salt")),
						},
					)
					counterApp.ConnectEndpoints(staticFS)
					counterApp.Logger.Println("Listening at", cCtx.String("bind"))
					counterApp.Serve()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
