package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func main() {
	m := martini.Classic()
	m.Get("/hello/:name", helloWorld)
	m.Post("/create", binding.Json(User{}), createUser)

	m.Run()
}

func helloWorld(params martini.Params) string {
	return "Hello, " + params["name"] + "!"
}

func createUser(user User) (int, string) {
	if len(user.Name) <= 0 || len(user.Pass) <= 0 {
		return 500, "User not created!"
	}

	// TODO: Save to the DB
	return 201, "User created successfully!"
}

type User struct {
	Name string `json:"name" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
