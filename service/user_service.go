package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Tags 用户
// @Summary 用户列表
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "",
		"data":    data,
	})
}

// GetUserByNameAndPassword
// @Tags 用户
// @Summary 根据用户名和密码查询用户
// @Param name postForm string false "用户名"
// @Param password postForm string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserByNameAndPassword [Post]
func GetUserByNameAndPassword(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	if name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名不能为空",
			"data":    "",
		})
		return
	}
	if password == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码不能为空",
			"data":    "",
		})
		return
	}
	user := models.GetUserByName(name)
	fmt.Println("hello")
	fmt.Println(user.ID)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户不存在或密码错误",
			"data":    "",
		})
		return
	}
	fmt.Println(password)
	if !utils.ValidPassword(password, user.Salt, user.Password) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户不存在或密码错误",
			"data":    "",
		})
		return
	}

	user = models.GetUserByNameAndPassword(user.Name, user.Password)
	data := models.ResponseData{
		ID:       user.ID,
		Name:     user.Name,
		Identity: user.Identity,
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "登陆成功",
		"data":    data,
	})
}

// CreateUser
// @Tags 用户
// @Summary 新增用户
// @Param name PostForm string false "用户名"
// @Param password PostForm string false "密码"
// @Param rePassword PostForm string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [Post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	rePassword := c.PostForm("rePassword")
	if user.Name == "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名不能为空，注册失败",
			"data":    "",
		})
		return
	}
	if password == "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "密码不能为空，注册失败",
			"data":    "",
		})
		return
	}
	if password != rePassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
			"data":    "",
		})
		return
	}
	user.Password = password
	data := models.GetUserByName(user.Name)
	if data.ID != 0 {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "创建用户失败，用户名已存在",
			"data":    "",
		})
		return
	}

	user.Salt = string(rand.Int31())
	user.Password = utils.MakePassword(user.Password, user.Salt)
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "创建用户成功",
		"data":    "",
	})
}

// DeleteUser
// @Tags 用户
// @Summary 删除用户
// @Param id query string false "用户Id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [Post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user = models.GetUser(uint(id))
	if user.ID == 0 {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "删除用户失败，对象不存在",
			"data":    "",
		})
		return
	}
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功",
		"data":    "",
	})
}

// UpdateUser
// @Tags 用户
// @Summary 更新用户
// @Param id formData string false "用户Id"
// @Param name formData string false "用户名"
// @Param password formData string false "密码"
// @Param phone formData string false "电话"
// @Param email formData string false "邮箱"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [Post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user = models.GetUser(uint(id))
	if user.ID != 0 {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("更新用户失败，对象不存在"),
			"data":    "",
		})
		return
	}
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "更新用户成功",
		"data":    "",
	})
}

// UpGrade 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}

}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(userId))
	userResponses := make([]models.ResponseData, 0)
	for _, v := range users {
		userResponses = append(userResponses, models.ResponseData{
			ID:       v.ID,
			Name:     v.Name,
			Identity: v.Identity,
		})
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "查询用户好友成功",
		"data":    userResponses,
	})
}
