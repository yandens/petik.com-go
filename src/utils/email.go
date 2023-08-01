package utils

import (
  "fmt"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "gopkg.in/gomail.v2"
)

func SendEmail(user models.User) {
  token, err := GenerateToken(user.ID, user.Email, user.Role.Role)
  if err != nil {
    fmt.Println("Could not generate token")
  }

  linkVerify := fmt.Sprintf("%s/auth/verify?token=%s", configs.GetEnv("APP_URL"), token)

  m := gomail.NewMessage()
  m.SetHeader("From", configs.GetEnv("EMAIL_USERNAME"))
  m.SetHeader("To", user.Email)
  m.SetHeader("Subject", "Verify Email")
  m.SetBody("text/html", fmt.Sprintf("Click this link to verify your email: <a href=\"%s\">Here</a>", linkVerify))

  d := gomail.NewDialer("smtp.gmail.com", 587, configs.GetEnv("EMAIL_USERNAME"), configs.GetEnv("EMAIL_PASSWORD"))

  // send email
  if err := d.DialAndSend(m); err != nil {
    fmt.Println("Could not send email because: ", err)
  }
}
