package main

import (
	"fmt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/database"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/server"
)

func main() {
	server := server.NewHttpServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	server.MountMiddlewares()
	server.MountRoutes(psqlDB)

	app := server.GetApp()
	routes := app.GetRoutes()

	// Log availables routes when initialized
	for _, route := range routes {
		fmt.Printf("%s -> '%s'\n", route.Method, route.Path)
	}

	server.Start(env.AppEnv.AppPort)
}
