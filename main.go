package main

import (
	"auth/handlers"
	"auth/utils"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main(){

	//loading the env files 
	utils.LoadEnv()

	//creating a new go server
	app := fiber.New()

	//connecting to the database 
	utils.ConnectDB()

	// specifying routes 
	app.Post("/register",handlers.Register)
	app.Post("/login",handlers.Login)

	//start the server
	app.Listen(os.Getenv("PORT"))
}