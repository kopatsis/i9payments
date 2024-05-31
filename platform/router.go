package platform

import (
	"i9pay/platform/login"
	"i9pay/platform/middleware"
	"i9pay/platform/multipass"
	"i9pay/platform/pay"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(auth *auth.Client, database *mongo.Database, scheduler *gocron.Scheduler) *gin.Engine {
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
	router.POST("/process-payment", pay.PostPayment(auth, database))
	router.POST("/cancel", pay.CancelPayment(auth, database, scheduler))
	router.POST("/update", pay.UpdateSubscriptionPaymentMethod())
	router.POST("/uncancel", pay.PostUncancel(auth, database, scheduler))
	router.POST("/swap", pay.UpdateFrequency(auth, database))

	return router
}
