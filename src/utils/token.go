package utils

import (
  "fmt"
  "github.com/golang-jwt/jwt/v5"
  "github.com/yandens/petik.com-go/src/configs"
  "time"
)

func GenerateToken(userID uint, email string, role string) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id":    userID,
    "email": email,
    "role":  role,
    "exp":   time.Now().Add(time.Hour * 24).Unix(),
  })

  tokenString, err := token.SignedString([]byte(configs.GetEnv("JWT_SECRET")))

  if err != nil {
    return "", err
  }

  return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }

    return []byte(configs.GetEnv("JWT_SECRET")), nil
  })

  // check if the token is valid
  if err != nil {
    return nil, err
  }

  return token, nil
}
