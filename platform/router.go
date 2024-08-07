package platform

import (
	"i9pay/platform/home"
	"i9pay/platform/login"
	"i9pay/platform/middleware"
	"i9pay/platform/multipass"
	"i9pay/platform/pay"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(auth *auth.Client, database *mongo.Database, client *sendgrid.Client) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.AuthMiddleware(auth))

	router.LoadHTMLGlob("html/*")
	router.Static("/static", "./static")

	router.GET("/multipass", multipass.Multipass(auth, database))
	// router.GET("/sub", pay.Subscription(auth, database))

	router.GET("/", home.AdminPanel(false, auth, database))
	router.GET("/mobile", home.AdminPanel(true, auth, database))

	router.GET("/code", multipass.SpecialCode(database))

	router.GET("/login", login.LoginPage(auth, database))
	router.GET("/new", login.SignUpPage(auth, database))
	router.GET("/logout", login.Logout(auth))
	router.GET("/pay", pay.Subscription(auth, database))

	router.POST("/updateName", home.Name(auth, database))
	router.POST("/delete", home.Delete(client, auth, database))

	router.POST("/verifyToken", login.VerifyToken(auth, database))
	router.POST("/process", pay.PostPayment(auth, database))
	router.POST("/cancel", pay.CancelPayment(client, auth, database))
	router.POST("/update", pay.UpdateSubscriptionPaymentMethod(auth, database))
	router.POST("/uncancel", pay.PostUncancel(client, auth, database))
	router.POST("/swap", pay.UpdateFrequency(auth, database))

	router.POST("/confirmationwh", pay.WebhookConfirm(client, auth, database))
	router.POST("/failedwh", pay.WebhookFail(client, auth, database))
	router.PATCH("/resetdate", login.ResetPasswordDate(auth, database))

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
	})

	return router
}
