package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/fzzy/radix/redis"
	"log"
	"io"
	"crypto/md5"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	errorHandler(err)
	defer conn.Close()
	m.Map(conn)

	// Endpoints
	m.Get("/hello/:name", helloWorld)
	m.Post("/create", binding.Json(User{}), createUser)
	m.Post("/authenticate", binding.Json(User{}), authenticateUser)

	// Start the server
	m.Run()
}

func helloWorld(params martini.Params) string {
	return "Hello, " + params["name"] + "!"
}

func createUser(user User, r render.Render, conn *redis.Client) {
	if (len(user.Name) <= 0 || len(user.Pass) <= 0 || hasUser(user, conn)) {
		r.JSON(500, "User not created.")
	} else {
		saveUser(user, conn)
		r.JSON(201, "User created successfully!")
	}
}

func authenticateUser(user User, r render.Render, conn *redis.Client) {

	fetchedUser := retrieveUser(user, conn)

	if (userToMD5(user) == fetchedUser.Pass) {
		r.JSON(200, "Welcome, "+user.Name+"!")
	} else {
		r.JSON(401, "User or password incorrect.")
	}
}

func hasUser(user User, conn *redis.Client) bool {
	return retrieveUser(user, conn).Pass != ""
}

func retrieveUser(user User, conn *redis.Client) User {
	fetchedPassword, err := conn.Cmd("get", user.Name).Str()  // TODO: Find better way to see if the Reply is null
	errorHandler(err)
	return User{user.Name, fetchedPassword}
}

func saveUser(user User, conn *redis.Client) {
	conn.Cmd("set", user.Name, userToMD5(user))
}

func userToMD5(user User) string {
	md5 := md5.New()
	io.WriteString(md5, user.Name)
	io.WriteString(md5, user.Pass)
	return string(md5.Sum(nil))
}

func errorHandler(error error) {
	if (error != nil) {
		log.Println("AN ERROR HAS OCCURED:")
		log.Println(error)
	}
}

type User struct {
	Name string `json:"name" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
