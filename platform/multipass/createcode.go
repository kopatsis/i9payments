package multipass

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SpecialCode(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, err := getSpecialCode(database)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Error generating code",
				"Exact": err,
			})
		}

		c.JSON(201, gin.H{"code": code})
	}
}
