package main

import (
	"context"

	application "github.com/zillalikestocode/community-api/api/app"
	"github.com/zillalikestocode/community-api/api/app/configs"
)

func main() {

	client := configs.ConnectDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := application.New()
	app.Start(context.TODO())
}
