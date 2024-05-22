package platform

import (
	"net/http"

	"i9pay/platform/login"
	"i9pay/platform/middleware"
	"i9pay/platform/multipass"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(auth *auth.Client, database *mongo.Database) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Static("/static", "./static")

	router.LoadHTMLGlob("../html/*")
	router.GET("/sub", multipass.Multipass(auth, database))

	router.GET("/code", multipass.SpecialCode(database))

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", nil)
	})

	router.POST("/verifyToken", login.VerifyToken(auth))

	return router
}
