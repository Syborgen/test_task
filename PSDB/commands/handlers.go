package commands

import (
	"PSBD/datastructures"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ShowObjectHandler(c *gin.Context) {
	objects, err := GetObjects()
	if err != nil {
		sendError("get objects error: "+err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, objects)
}

func ShowWindowHandler(c *gin.Context) {
	techWindows, err := GetTechWindows()
	if err != nil {
		sendError("get tech windows error: "+err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, techWindows)
}

func ShowWindowAllHandler(c *gin.Context) {
	techWindows, err := GetTechWindowsAll()
	if err != nil {
		sendError("get tech windows all error: "+err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, techWindows)
}

func ShowWindowSortHandler(c *gin.Context) {
	sort := c.Query("sort")
	start := c.Query("start")
	end := c.Query("end")

	action := c.Query("action")
	var techWindows []datastructures.GroupedTechWindow
	var err error
	switch action {
	case "query":
		techWindows, err = GetSortedWindowsQuery(sort, start, end)
		if err != nil {
			sendError("get sorted windows query error: "+err.Error(), c)
			return
		}

	case "proc":
		techWindows, err = GetSortedWindowsProc(sort, start, end)
		if err != nil {
			sendError("get sorted windows proc error: "+err.Error(), c)
			return
		}

	}

	c.JSON(http.StatusOK, techWindows)
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
		sendError("bind json error: "+err.Error(), c)
		return
	}

	switch addWindowArguments.Action {
	case "query":
		err = AddWindowQuery(
			addWindowArguments.ObjectID,
			addWindowArguments.Start,
			addWindowArguments.End,
		)
		if err != nil {
			sendError("add window query error: "+err.Error(), c)
			return
		}

	case "proc":
		err = AddWindowProc(
			addWindowArguments.ObjectID,
			addWindowArguments.Start,
			addWindowArguments.End,
		)
		if err != nil {
			sendError("add window proc error: "+err.Error(), c)
			return
		}

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
		sendError("bind json error: "+err.Error(), c)
		return
	}

	err = Create(createArguments.Objects, createArguments.Windows)
	if err != nil {
		sendError("create command error: "+err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, ServerResponse{Status: "ok"})
}

func sendError(errorText string, c *gin.Context) {
	fmt.Println(errorText)
	response := ServerResponse{Status: "error", Message: errorText}
	c.JSON(http.StatusOK, response)
}
