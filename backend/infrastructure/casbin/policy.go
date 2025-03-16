// infrastructure/casbin/policy.go
package casbin

import (
	"github.com/casbin/casbin/v2"
)

// CasbinPolicyChecker ใช้ PolicyChecker ด้วย Casbin
type CasbinPolicyChecker struct {
	enforcer *casbin.Enforcer
}

// NewCasbinPolicyChecker สร้าง CasbinPolicyChecker ใหม่
func NewCasbinPolicyChecker(modelPath, policyPath string) (*CasbinPolicyChecker, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	return &CasbinPolicyChecker{
		enforcer: enforcer,
	}, nil
}

// CheckPolicy ตรวจสอบว่าบทบาทมีสิทธิ์ในการดำเนินการกับทรัพยากรหรือไม่
func (c *CasbinPolicyChecker) CheckPolicy(role, resource, action string) (bool, error) {
	return c.enforcer.Enforce(role, resource, action)
}
