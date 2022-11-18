package main

import (
	"PSBD/commands"
	"PSBD/dbhelper"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// sudo docker build --tag docker-psdb .
func main() {
	var err error
	dbhelper.DB, err = dbhelper.ConnectToDb()
	if err != nil {
		fmt.Println("Connect to db error:", err)
		return
	}

	defer dbhelper.DB.Close()

	router := gin.Default()

	router.GET("/show_object", commands.ShowObjectHandler)
	router.GET("/show_window", commands.ShowWindowHandler)
	router.GET("/show_window_sort", commands.ShowWindowSortHandler)
	router.GET("/show_window_all", commands.ShowWindowAllHandler)
	router.POST("/add_window", commands.AddWindowHandler)
	router.POST("/create", commands.CreateHandler)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Server start error:", err)
	}
}
