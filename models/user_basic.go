package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserBasic struct {
	gorm.Model
	ID            uint   `json:"id"`
	Identity      string `json:"identity"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$))"`
	Email         string `valid:"email"`
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogOut      bool
	DeviceInfo    string
}

type ResponseData struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Identity string `json:"identity"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUser(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).Find(&user)
	return user
}

func GetUserList() []*UserBasic {
	user := make([]*UserBasic, 10)
	utils.DB.Find(&user)
	return user
}

func GetUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func GetUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func GetUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func GetUserByNameAndPassword(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? AND password = ?", name, password).First(&user)
	// token encode
	str := strconv.FormatInt(time.Now().Unix(), 10)
	temp := utils.Md5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

// DeleteUser 逻辑删除
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}
