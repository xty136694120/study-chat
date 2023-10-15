package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigFile("config/app.yml")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Println("config app init")

}

func InitMySQL() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true},
	)
	var err error
	// sudo /usr/local/mysql/support-files/mysql.server start
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("config MySQL init")
}

func InitRedis() {
	// bin: redis-server
	var ctx = context.Background()
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := Redis.Ping(ctx).Result()
	if err != nil {
		fmt.Println("failed to connect redis: ", err)
	} else {
		fmt.Println("config Redis init: ", pong)
	}
}

const (
	PublishKey = "websocket"
)

func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	err = Redis.Publish(ctx, channel, msg).Err()
	return err
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Redis.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	return msg.Payload, err
}
