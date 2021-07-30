package main

import (
	"casbin-gorm/middleware"
	"casbin-gorm/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gormAdapter "github.com/casbin/gorm-adapter/v3"
)

var authorities = []model.Authority{
	{ID: "1", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "admin"},
	{ID: "2", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "visitor"},
}

var users = []model.User{
	{ID: "1", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "admin", Email: "admin@163.com", Password: "123456.", AuthorityID: "1", Mobile: "15288888888"},
	{ID: "2", CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "visitor", Email: "visitor@163.com", Password: "123456...", AuthorityID: "2", Mobile: "15766666666"},
}

var carbinRules = []gormAdapter.CasbinRule{
	{ID: 1, Ptype: "p", V0: "1", V1: "/users", V2: "GET"},
	{ID: 2, Ptype: "p", V0: "1", V1: "/user/:id", V2: "GET"},
	{ID: 3, Ptype: "p", V0: "1", V1: "/user", V2: "PUT"},
	{ID: 4, Ptype: "p", V0: "1", V1: "/user", V2: "POST"},
	{ID: 5, Ptype: "p", V0: "1", V1: "/user/:id", V2: "DELETE"},
	{ID: 6, Ptype: "p", V0: "1", V1: "/user/:id/tag", V2: "GET"},
	{ID: 7, Ptype: "p", V0: "2", V1: "/user", V2: "POST"},
	{ID: 8, Ptype: "p", V0: "2", V1: "/user/:id", V2: "DELETE"},
}

func main() {
	dsn := "root:0000@tcp(127.0.0.1:3306)/gorm1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&model.User{},
		&model.Authority{},
	)

	if !db.Migrator().HasTable("casbin_rule") {
		db.Migrator().CreateTable(&gormAdapter.CasbinRule{})
	}

	initData(db)

	router := gin.Default()

	router.Use(middleware.CasbinHandler())

	router.GET("/users", func(c *gin.Context) {
		c.String(http.StatusOK, "get user list")
	})

	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "get user by id: %s", id)
	})

	router.GET("/user/:id/tag", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "get user tag by id: %s", id)
	})

	router.POST("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "create user")
	})

	router.PUT("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "update user")
	})

	router.DELETE("/user/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "delete user by id")
	})

	router.Run(":7788")
}

func initData(db *gorm.DB) {
	db.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []string{"1", "2"}).Find(&[]model.Authority{}).RowsAffected == 2 {
			fmt.Println("数据已经初始化")
			return nil
		}
		if err := tx.Create(&authorities).Error; err != nil {
			return err
		}
		return nil
	})

	db.Transaction(func(tx *gorm.DB) error {
		if tx.Find(&[]gormAdapter.CasbinRule{}).RowsAffected == 8 {
			fmt.Println("数据已经初始化")
			return nil
		}
		if err := tx.Create(&carbinRules).Error; err != nil {
			return err
		}

		return nil
	})

	db.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []string{"1", "2"}).Find(&[]model.User{}).RowsAffected == 2 {
			fmt.Println("数据已经初始化")
			return nil
		}
		if err := tx.Create(&users).Error; err != nil {
			return err
		}
		return nil
	})
}
