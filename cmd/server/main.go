package main

import (
	"context"
	"github.com/zdarovich/win-lose-api/api"
	"github.com/zdarovich/win-lose-api/cronjobs"
	"github.com/zdarovich/win-lose-api/database"
)

func main() {
	database.InitDB()
	route := api.NewRouter()
	apiEngine := api.New(route)
	cronjobs.GetInstance().Run(context.Background())
	apiEngine.Run(8081)
}