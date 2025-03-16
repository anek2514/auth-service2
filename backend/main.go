package main

import (
	"auth-service/config"
	"auth-service/infrastructure/casbin"
	"auth-service/infrastructure/keycloak"
	"auth-service/interface/api"
	"auth-service/interface/middleware"
	"auth-service/usecase"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// โหลดการตั้งค่า
	cfg := config.LoadConfig()

	log.Printf("เริ่มต้นบริการยืนยันตัวตนบน port %s", cfg.Port)
	log.Printf("กำหนดค่า KeycloakIssuer: %s", cfg.KeycloakIssuer)
	log.Printf("กำหนดค่า ModelPath: %s", cfg.ModelPath)
	log.Printf("กำหนดค่า PolicyPath: %s", cfg.PolicyPath)

	// เริ่มต้นตัวตรวจสอบโทเค็น Keycloak
	tokenVerifier, err := keycloak.NewKeycloakTokenVerifier(cfg.KeycloakPublicKey, cfg.KeycloakIssuer)
	if err != nil {
		log.Fatalf("ไม่สามารถเริ่มต้นตัวตรวจสอบโทเค็น Keycloak: %v", err)
	}

	// เริ่มต้นตัวตรวจสอบนโยบาย Casbin
	policyChecker, err := casbin.NewCasbinPolicyChecker(cfg.ModelPath, cfg.PolicyPath)
	if err != nil {
		log.Fatalf("ไม่สามารถเริ่มต้นตัวตรวจสอบนโยบาย Casbin: %v", err)
	}

	// เริ่มต้นบริการยืนยันตัวตน (use case)
	authService := usecase.NewAuthService(tokenVerifier, policyChecker)

	// เริ่มต้นเราเตอร์ Gin
	r := gin.Default()

	// ตั้งค่า CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ตรวจสอบสถานะ
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/data", api.DataHandler)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("ไม่สามารถเริ่มต้นเซิร์ฟเวอร์: %v", err)
	}
}