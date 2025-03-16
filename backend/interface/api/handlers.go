// interface/api/handlers.go
package api

import (
	"github.com/gin-gonic/gin"
)

// DataHandler จัดการคำขอสำหรับข้อมูลที่ป้องกัน
func DataHandler(c *gin.Context) {
	username := c.GetString("username")
	roles := c.GetStringSlice("roles")

	c.JSON(200, gin.H{
		"data":     "นี่คือข้อมูลที่ได้รับการป้องกัน",
		"username": username,
		"roles":    roles,
	})
}
