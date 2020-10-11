package auth

// Auth 加解密相关
type Auth struct {
	Key string
}

// NewAuth 工厂方法
func NewAuth(k string) *Auth {
	return &Auth{
		Key: k,
	}
}
