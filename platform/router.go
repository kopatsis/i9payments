package platform

import (
	"i9pay/platform/login"
	"i9pay/platform/middleware"
	"i9pay/platform/multipass"
	"i9pay/platform/pay"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(auth *auth.Client, database *mongo.Database) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.AuthMiddleware(auth))

	router.Static("/static", "./static")

	router.LoadHTMLGlob("../html/*")
	router.GET("/multipass", multipass.Multipass(auth, database))
	router.GET("/sub", pay.Subscription(auth, database))

	router.GET("/code", multipass.SpecialCode(database))

	router.GET("/login", login.LoginPage(auth))
	router.GET("/new", login.SignUpPage(auth))

	router.POST("/verifyToken", login.VerifyToken(auth))
	router.POST("/process-payment", pay.PostPayment(auth))
	router.POST("/cancel", pay.CancelPayment(auth, database))
	router.POST("/update", pay.UpdateSubscriptionPaymentMethod())

	return router
}
