package main

import (
	"issue-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/issue", handlers.CreateIssue)
	r.GET("/issues", handlers.ListIssues)
	r.GET("/issue/:id", handlers.GetIssueByID)
	r.PATCH("/issue/:id", handlers.UpdateIssue)

	r.Run(":8080")
}
