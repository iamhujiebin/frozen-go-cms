package jwt

import (
	"frozen-go-cms/common/resource/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 载荷，增加用户别名
type Claims struct {
	UserId uint64
	Mobile string
	jwt.StandardClaims
}

// 生成token
func GenerateToken(userId uint64, mobile, issuer string) (string, error) {
	jwtConfig := config.GetConfigJWT()
	duration, err := time.ParseDuration(jwtConfig.EXPIRE)
	if err != nil {
		return "", err
	}

	expireTime := time.Now().Add(duration)
	claims := Claims{
		UserId: userId,
		Mobile: mobile,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			Issuer:    issuer,            //签名的发行者
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func GetJWTSecret() []byte {
	return []byte(config.GetConfigJWT().SECRET)
}

// 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
