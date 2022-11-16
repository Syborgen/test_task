package commands

import (
	"PSBD/dbhelper"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowObjectHandler(c *gin.Context) {
	objects := GetObjects()
	c.JSON(http.StatusOK, objects)
}

func ShowWindowHandler(c *gin.Context) {
	techWindows := GetTechWindows()
	c.JSON(http.StatusOK, techWindows)
}

func ShowWindowSortHandler(c *gin.Context) {
	sort := c.Query("sort")
	start := c.Query("start")
	end := c.Query("end")

	action := c.Query("action")
	var techWindows []dbhelper.GroupedTechWindow
	switch action {
	case "query":
		techWindows = GetSortedWindowsQuery(sort, start, end)
	case "proc":
		techWindows = GetSortedWindowsProc(sort, start, end)
	}

	c.JSON(http.StatusOK, techWindows)
}

type ServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AddWindowCommandArguments struct {
	ObjectID int    `json:"object_id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Action   string `json:"action"`
}

func AddWindowHandler(c *gin.Context) {
	var addWindowArguments AddWindowCommandArguments
	err := c.BindJSON(&addWindowArguments)
	if err != nil {
		fmt.Println("Bind JSON error:", err)
	}

	switch addWindowArguments.Action {
	case "query":
		err = AddWindowQuery(
			addWindowArguments.ObjectID,
			addWindowArguments.Start,
			addWindowArguments.End,
		)
	case "proc":
		err = AddWindowProc(
			addWindowArguments.ObjectID,
			addWindowArguments.Start,
			addWindowArguments.End,
		)
	}

	if err != nil {
		response := ServerResponse{Status: "error", Message: err.Error()}
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, ServerResponse{Status: "ok"})
}

type CreateCommandArguments struct {
	Objects int `json:"objects"`
	Windows int `json:"windows"`
}

func CreateHandler(c *gin.Context) {
	var createArguments CreateCommandArguments
	err := c.BindJSON(&createArguments)
	if err != nil {
		fmt.Println("Bind JSON error:", err)
	}

	Create(createArguments.Objects, createArguments.Windows)
	c.JSON(http.StatusOK, nil)
}
