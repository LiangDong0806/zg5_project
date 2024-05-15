package models

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"strconv"
	"zg5/work01/server/common/global"
)

func InitMysql(mysqlFunc func(DB *gorm.DB) error) error {
	username := global.ServerConfig.MySQL.Username
	password := global.ServerConfig.MySQL.Password
	host := global.ServerConfig.MySQL.Host
	port := strconv.Itoa(global.ServerConfig.MySQL.Port)
	dbname := global.ServerConfig.MySQL.Library
	//dsn := "root:123456@tcp(127.0.0.1:3307)/JobDirectory?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println(dsn, "22222222222222222222222222222")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	exit, _ := db.DB()
	defer exit.Close()
	err = mysqlFunc(db)
	return err
}
func Migrate() {
	err := InitMysql(func(DB *gorm.DB) error {
		return DB.AutoMigrate(&User{})
	})
	if err != nil {
		panic(err)
	}
}

func GetUserByUsername(username string) (*User, error) {
	user := new(User)
	return user, InitMysql(func(DB *gorm.DB) error {
		res := DB.Where("username = ?", username).First(&user)
		if res.RowsAffected == 0 {
			return errors.New("账号不存在")
		}
		return nil
	})
}

type User struct {
	gorm.Model
	Username string
	Password string
}
