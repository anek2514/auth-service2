// infrastructure/keycloak/keycloak.go
package keycloak

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"auth-service/domain"

	"github.com/golang-jwt/jwt/v4"
)

// KeycloakTokenVerifier ใช้ TokenVerifier สำหรับโทเค็น Keycloak
type KeycloakTokenVerifier struct {
	publicKey *rsa.PublicKey
	issuer    string
}

// NewKeycloakTokenVerifier สร้าง KeycloakTokenVerifier ใหม่
func NewKeycloakTokenVerifier(publicKeyPEM string, issuer string) (*KeycloakTokenVerifier, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("ไม่สามารถแยกบล็อก PEM ที่มีคีย์สาธารณะ")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถแยกคีย์สาธารณะที่เข้ารหัส DER: %v", err)
	}

	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("คีย์ไม่ใช่ประเภท *rsa.PublicKey")
	}

	return &KeycloakTokenVerifier{
		publicKey: rsaPublicKey,
		issuer:    issuer,
	}, nil
}

// VerifyToken ตรวจสอบโทเค็น JWT ของ Keycloak และดึงข้อมูล
func (v *KeycloakTokenVerifier) VerifyToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("วิธีการลงนามไม่ตรงกับที่คาดหวัง: %v", token.Header["alg"])
		}
		return v.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("โทเค็นไม่ถูกต้อง: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("ข้อมูลในโทเค็นไม่ถูกต้อง")
	}

	// ตรวจสอบวันหมดอายุของโทเค็น
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("โทเค็นหมดอายุแล้ว")
	}

	// ตรวจสอบผู้ออกโทเค็น
	if !claims.VerifyIssuer(v.issuer, true) {
		return nil, errors.New("ผู้ออกโทเค็นไม่ถูกต้อง")
	}

	// ดึงชื่อผู้ใช้
	username, ok := claims["preferred_username"].(string)
	if !ok {
		return nil, errors.New("ไม่พบชื่อผู้ใช้ในโทเค็น")
	}

	// ดึงบทบาท
	realmAccess, ok := claims["realm_access"].(map[string]interface{})
	if !ok {
		return nil, errors.New("ไม่พบบทบาทในโทเค็น")
	}

	rawRoles, ok := realmAccess["roles"].([]interface{})
	if !ok || len(rawRoles) == 0 {
		return nil, errors.New("ไม่พบบทบาทในโทเค็น")
	}

	var roles []string
	for _, r := range rawRoles {
		if roleStr, ok := r.(string); ok {
			roles = append(roles, roleStr)
		}
	}

	return &domain.Claims{
		Username: username,
		Roles:    roles,
	}, nil
}
