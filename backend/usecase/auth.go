// usecase/auth.go
package usecase

import (
	"auth-service/domain"
	"errors"
)

// AuthUseCase กำหนดอินเทอร์เฟซของ use case การยืนยันตัวตน
type AuthUseCase interface {
	VerifyAndExtractClaims(token string) (*domain.Claims, error)
	CheckPermission(roles []string, resource string, action string) (bool, error)
}

// AuthService ทำงานตามอินเทอร์เฟซ AuthUseCase
type AuthService struct {
	tokenVerifier TokenVerifier
	policyChecker PolicyChecker
}

// อินเทอร์เฟซช่วยให้เราสามารถเปลี่ยนการใช้งานได้โดยไม่ต้องเปลี่ยน Busniss Logic
type TokenVerifier interface {
	VerifyToken(tokenString string) (*domain.Claims, error)
}

type PolicyChecker interface {
	CheckPolicy(role string, resource string, action string) (bool, error)
}

// NewAuthService สร้าง AuthService ใหม่
func NewAuthService(tokenVerifier TokenVerifier, policyChecker PolicyChecker) *AuthService {
	return &AuthService{
		tokenVerifier: tokenVerifier,
		policyChecker: policyChecker,
	}
}

// VerifyAndExtractClaims ตรวจสอบโทเค็น JWT และดึงข้อมูลที่เกี่ยวข้อง
func (s *AuthService) VerifyAndExtractClaims(token string) (*domain.Claims, error) {
	if token == "" {
		return nil, errors.New("token ว่างเปล่า")
	}

	return s.tokenVerifier.VerifyToken(token)
}

// CheckPermission ตรวจสอบว่ามีบทบาทใดบ้างที่มีสิทธิ์เข้าถึงทรัพยากร
func (s *AuthService) CheckPermission(roles []string, resource string, action string) (bool, error) {
	for _, role := range roles {
		allowed, err := s.policyChecker.CheckPolicy(role, resource, action)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}
