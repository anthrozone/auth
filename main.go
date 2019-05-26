package main

import (
	"github.com/globalsign/mgo"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/auth/platform"
	"log"
	"os"
)

func main() {
	e := echo.New()
	blog := new(platform.Platform)

	blog.Key = "mysuperawesometestkey"

	var err error

	blog.Mongo, err = mgo.Dial(os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	e.POST("/auth/login", blog.Login)
	e.POST("/auth/signup", blog.Register)

	e.Logger.Fatal(e.Start(":8080"))
}
