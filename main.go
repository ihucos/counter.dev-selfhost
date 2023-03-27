package main

import (
	rand "crypto/rand"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	"github.com/gomodule/redigo/redis"
	_ "github.com/ihucos/counter.dev/endpoints"
	"github.com/ihucos/counter.dev/lib"
	"github.com/ihucos/counter.dev/models"
	"golang.org/x/term"

	"github.com/urfave/cli/v2"
)

//go:embed all:static
var staticFS embed.FS

func getSecret(conn redis.Conn, key string) string {

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

func getApp(cCtx *cli.Context) *lib.App {
	conn, err := redis.DialURL(cCtx.String("redis-url"))
	if err != nil {
		fmt.Printf("Can't connect to redis server at %s - %s\n", cCtx.String("redis-url"), err)
		os.Exit(1)
	}
	defer conn.Close()
	return lib.NewApp(
		lib.Config{
			RedisUrl:     cCtx.String("redis-url"),
			Bind:         cCtx.String("bind"),
			CookieSecret: []byte(getSecret(conn, "cntr:config:cookie_secret")),
			PasswordSalt: []byte(getSecret(conn, "cntr:config:password_salt")),
		},
	)

}

func NewUser(app *lib.App, userID string) models.User {
	conn := app.RedisPool.Get()
	return models.NewUser(conn, userID, app.Config.PasswordSalt)
}

func main() {

	redis_flag := &cli.StringFlag{
		Name:  "redis-url",
		Value: "redis://localhost:6379",
		Usage: "Which redis server to connect to",
	}

	cliApp := &cli.App{
		Flags: []cli.Flag{redis_flag},
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
					redis_flag,
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

					app := getApp(cCtx)
					app.ConnectEndpoints(staticFS)
					app.Logger.Println("Listening at", cCtx.String("bind"))
					app.Serve()
					return nil
				},
			},
			{
				Name:      "createuser",
				Usage:     "Create a new user",
				ArgsUsage: "<username>",
				Flags: []cli.Flag{redis_flag,
					&cli.IntFlag{
						Required: true,
						Name:  "utc-offset",
						Usage: "Specify your favorite timezones utc offset",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() < 1 {
						fmt.Println("Missing argument: user")
						os.Exit(1)
					}
					app := getApp(cCtx)
					user := NewUser(app, cCtx.Args().Get(0))
					fmt.Print("Password for new user: ")
					password, err := term.ReadPassword(0)
					fmt.Print("\n")
					if err != nil {
						return err
					}
					err = user.Create(string(password))
					if err != nil {
						return err
					}
					return user.SetPref("utcoffset", cCtx.String("utc-offset"))
				},
			},
			{
				Name:      "chgpwd",
				Usage:     "Change a users password",
				ArgsUsage: "<username>",
				Flags:     []cli.Flag{redis_flag},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() < 1 {
						fmt.Println("Missing argument: user")
						os.Exit(1)
					}
					app := getApp(cCtx)
					userID := cCtx.Args().Get(0)
					user := NewUser(app, userID)
					fmt.Printf("New password user %s: ", userID)
					newPassword, err := term.ReadPassword(0)
					fmt.Print("\n")
					if err != nil {
						return err
					}
					return user.ChangePassword(string(newPassword))
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)

	}
}
