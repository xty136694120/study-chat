package main

import (
	"ginchat/models"
	"ginchat/utils"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	utils.InitConfig()
	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	//err = db.AutoMigrate(&models.UserBasic{})
	err = db.AutoMigrate(&models.Message{})
	err = db.AutoMigrate(&models.Contact{})
	err = db.AutoMigrate(&models.GroupBasic{})
	//if err != nil {
	//	return
	//}

	// Create
	//user := &models.UserBasic{}
	//user.Name = "xty"
	//user.Password = "123"
	//db.Create(user)

	// Read
	//fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//db.First(user, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	//db.Model(user).Update("Password", "1234")
	// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	//db.Delete(&product, 1)
}
