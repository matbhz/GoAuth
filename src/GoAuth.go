package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/fzzy/radix/redis"
	"log"
)

func errorHandler(err error) {
	log.Println(err)
}


func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	errorHandler(err)
	defer conn.Close()

	m.Map(conn)

	m.Get("/hello/:name", helloWorld)
	m.Post("/create", binding.Json(User{}), createUser)

	m.Run()
}

func helloWorld(params martini.Params) string {
	return "Hello, " + params["name"] + "!"
}

func createUser(user User, r render.Render, conn *redis.Client) {
	if len(user.Name) <= 0 || len(user.Pass) <= 0 || hasUser(user, conn) {
		r.JSON(500, "User not created!")
	} else {
		conn.Cmd("set", user.Name, user.Pass)
		r.JSON(201, "User created successfully!")
	}
}

func hasUser(user User, conn *redis.Client) bool {
	value, err := conn.Cmd("get", user.Name).Str() // TODO: Find better way to see if the Reply is null
	errorHandler(err)
	return value != ""
}

type User struct {
	Name string `json:"name" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
