package utils

import (
  "fmt"
  "github.com/golang-jwt/jwt/v5"
  "github.com/yandens/petik.com-go/src/configs"
  "strconv"
  "time"
)

func GenerateToken(userID uint, email string, role string) (string, error) {
  encryptedID, err := Encrypt([]byte(strconv.Itoa(int(userID))))
  if err != nil {
    return "", err
  }

  encryptedEmail, err := Encrypt([]byte(email))
  if err != nil {
    return "", err
  }

  encryptedRole, err := Encrypt([]byte(role))
  if err != nil {
    return "", err
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id":    encryptedID,
    "email": encryptedEmail,
    "role":  encryptedRole,
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

func TokenClaims(token *jwt.Token) (jwt.MapClaims, error) {
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return nil, fmt.Errorf("Invalid token")
  }

  // decrypt user id
  decryptedUserID, err := Decrypt([]byte(claims["id"].(string)))
  if err != nil {
    return nil, err
  }

  // decrypt user email
  decryptedEmail, err := Decrypt([]byte(claims["email"].(string)))
  if err != nil {
    return nil, err
  }

  // decrypt user role
  decryptedRole, err := Decrypt([]byte(claims["role"].(string)))
  if err != nil {
    return nil, err
  }

  // set claims
  claims["id"] = decryptedUserID
  claims["email"] = decryptedEmail
  claims["role"] = decryptedRole

  return claims, nil
}
