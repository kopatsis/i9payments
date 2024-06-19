package multipass

import (
	"fmt"
	"i9pay/platform/login"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Multipass(authClient *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		refresh := c.Query("multipass")
		code := c.Query("code")
		destination := c.Query("dest")

		if status := checkSpecialCode(code, database); !status || refresh == "" {
			fmt.Println("failed on code")
			c.Redirect(http.StatusFound, "/pay")
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		token, refresh, err := login.GetNewIDToken(refresh)
		if err != nil {
			fmt.Println("failed on get tokens")
			c.Redirect(http.StatusFound, "/pay")
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		fmt.Println(refresh)
		fmt.Println(token)

		if err := login.Cookie(token, refresh, authClient, c); err != nil {
			fmt.Println("failed on cookie create")
			c.Redirect(http.StatusFound, "/pay")
		}

		go func() {
			deleteSpecialCode(code, database)
		}()

		if destination == "pay" {
			c.Redirect(http.StatusFound, "/pay")
		} else if destination == "mobile" {
			c.Redirect(http.StatusFound, "/mobile")
		} else {
			c.Redirect(http.StatusFound, "/")
		}

	}
}
