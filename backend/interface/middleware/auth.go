// interface/middleware/auth.go
package middleware

import (
	"strings"

	"auth-service/usecase"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware สร้าง middleware ของ Gin สำหรับการยืนยันตัวตนและการอนุญาตด้วย JWT
func AuthMiddleware(authUseCase usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "ไม่ได้รับอนุญาต"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจสอบโทเค็นและดึงข้อมูล
		claims, err := authUseCase.VerifyAndExtractClaims(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// ตรวจสอบสิทธิ์สำหรับทรัพยากรที่ร้องขอ
		resource := c.Request.URL.Path
		action := c.Request.Method
		allowed, err := authUseCase.CheckPermission(claims.Roles, resource, action)
		if err != nil {
			c.JSON(500, gin.H{"error": "เกิดข้อผิดพลาดในการตรวจสอบสิทธิ์"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(403, gin.H{"error": "ไม่มีสิทธิ์เข้าถึง: สิทธิ์ไม่เพียงพอ"})
			c.Abort()
			return
		}

		// เก็บข้อมูลสำหรับ handler
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}
