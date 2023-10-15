package service

import (
	"fmt"
	"ginchat/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
	}
}

func Register(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}
	err = ind.Execute(c.Writer, "register")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
	}
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles(
		"views/chat/index.html",
		"views/chat/main.html",
		"views/chat/foot.html",
		"views/chat/userinfo.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/group.html",
		"views/chat/concat.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		return
	}
	user := models.UserBasic{}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user.ID = uint(userId)
	user.Identity = token
	err = ind.Execute(c.Writer, user)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
	}
}
