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
		tokenID := c.Query("multipass")
		code := c.Query("code")
		destination := c.Query("dest")

		refresh, err := idToRefreshToken(tokenID, database)
		if err != nil {
			fmt.Println("failed on get tokens")
			setRedirect(destination, c)
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		if status := checkSpecialCode(code, database); !status || refresh == "" {
			fmt.Println("failed on code")
			setRedirect(destination, c)
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		token, refresh, err := login.GetNewIDToken(refresh)
		if err != nil {
			fmt.Println("failed on get tokens")
			setRedirect(destination, c)
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		if err := login.Cookie(token, refresh, authClient, c); err != nil {
			fmt.Println("failed on cookie create")
			setRedirect(destination, c)
		}

		go func() {
			deleteSpecialCode(code, database)
		}()

		setRedirect(destination, c)
	}
}

func setRedirect(destination string, c *gin.Context) {
	if destination == "pay" {
		c.Redirect(http.StatusFound, "/pay")
	} else if destination == "mobile" {
		c.Redirect(http.StatusFound, "/mobile")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
