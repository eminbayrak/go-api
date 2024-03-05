// File: internal/handlers/admin_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func AdminHandler(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)

	var json struct {
		Value string `json:"value" binding:"required"`
	}

	if c.Bind(&json) == nil {
		db[user] = json.Value
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
