package main

import (
	"fmt"
	"http-server/cmd/server/config"
	"http-server/infras/db"
	"http-server/internal/server"
)

func main() {
	configInstance, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot load env")
	}

	dsn := config.GetDsn(&configInstance)
	DB, err := db.NewDB(dsn)
	if err != nil {
		panic("Cannot connect to database")
	}

	svr := server.NewServer(DB)
	svr.Start()
}
