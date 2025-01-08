package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

var db *gorm.DB

func main() {
	var err error
	dsn := "host=postgres user=postgres password=postgres dbname=users port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user User
		// SQLインジェクションの脆弱性を持つクエリ
		db.Raw("SELECT * FROM users WHERE username = '" + username + "' AND password = '" + password + "'").Scan(&user)

		if user.ID != 0 {
			c.HTML(http.StatusOK, "success.html", gin.H{"username": user.Username})
		} else {
			c.HTML(http.StatusUnauthorized, "failure.html", nil)
		}
	})

	r.Static("/static", "./static")
	r.LoadHTMLFiles("templates/login.html", "templates/success.html", "templates/failure.html")
	r.Run(":80")
}
