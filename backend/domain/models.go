// domain/models.go
package domain

// User แทนเอนทิตี้ผู้ใช้ในระบบของเรา
type User struct {
	Username string
	Roles    []string
}

// Claims แทนข้อมูลใน JWT ที่เราดึงออกมาเพื่อการอนุญาต
type Claims struct {
	Username string
	Roles    []string
}
