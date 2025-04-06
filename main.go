package main

import (
	"go_auth/controllers"
	"go_auth/initializers"
	"go_auth/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
	initializers.InitCache()
}

func main() {
	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/users", controllers.GetAllUsers)
	r.GET("/validate", middlewares.RequireAuth, controllers.Validate)

	r.GET("/users/:id", controllers.GetUser)
	r.PUT("/users/:id", controllers.PutUser)
	r.PATCH("/users/:id", controllers.PatchUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/reset_password", controllers.ResetPassword)

	r.Run(os.Getenv("ADDR"))
}
