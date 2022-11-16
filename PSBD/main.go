package main

import (
	"PSBD/commands"
	"PSBD/dbhelper"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// sudo docker build --tag docker-psdb .
func main() {
	dbhelper.DB = dbhelper.ConnectToDb()
	defer dbhelper.DB.Close()

	router := gin.Default()

	router.GET("/show_object", commands.ShowObjectHandler)
	router.GET("/show_window", commands.ShowWindowHandler)
	router.GET("/show_window_sort", commands.ShowWindowSortHandler)
	router.POST("/add_window", commands.AddWindowHandler)
	router.POST("/create", commands.CreateHandler)

	router.Run("0.0.0.0:8080")
}
