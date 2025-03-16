package config

import (
	"os"
	"path/filepath"
)

// Config เก็บการตั้งค่าแอปพลิเคชัน
type Config struct {
	KeycloakPublicKey string
	KeycloakIssuer    string
	ModelPath         string
	PolicyPath        string
	Port              string
	AllowedOrigins    []string
}

// LoadConfig โหลดการตั้งค่าจากตัวแปรสภาพแวดล้อมพร้อมค่าเริ่มต้น
func LoadConfig() *Config {
	allowedOrigins := []string{"http://localhost:3000"}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = []string{origins}
	}

	// ค่าเริ่มต้นสำหรับ public key (สำหรับการพัฒนา)
	defaultPublicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhVpkIC+zXBM4sFug3e7fMKYX8b+zXEJvPe1fRvthnG2a7xzOaIl8OXwgO3ldpFR576cAP1vNMEhI97kZ6vovm2YYtvm/yWJgKLPmEoRa3hxtoilcLA4PfPFJvYuf09RCHIFlzSXa2Up3h2uywLSPPAwkELfQKL1sMZitCfxp/JUXxCF5Kdjqg9EEgeXVXcBfizApUm1okZJKT+n7EQ3Ys0XdjkPgiEUUjw6SqrtEntcbniU6BfNjl/8GT/n9AThJDNlB+cgusXoA9WSdik0x8Mx8JNPDwWvVA0yaYBsTmHM1u1G+L9+G/lhtUnIWKt/0zilPmuneQyI+F5PmYQA7RwIDAQAB
-----END PUBLIC KEY-----`

	// กำหนดโฟลเดอร์สำหรับไฟล์ config
	configDir := "./config"
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		configDir = dir
	}

	return &Config{
		KeycloakPublicKey: getEnv("KEYCLOAK_PUBLIC_KEY", defaultPublicKey),
		KeycloakIssuer:    getEnv("KEYCLOAK_ISSUER", "http://keycloak:8080/realms/auth101"),
		ModelPath:         filepath.Join(configDir, getEnv("MODEL_FILENAME", "model.conf")),
		PolicyPath:        filepath.Join(configDir, getEnv("POLICY_FILENAME", "policy.csv")),
		Port:              getEnv("PORT", "8081"),
		AllowedOrigins:    allowedOrigins,
	}
}

// getEnv รับตัวแปรสภาพแวดล้อมหรือคืนค่าเริ่มต้น
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}