package helpers

import (
  "bytes"
  "fmt"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "gopkg.in/gomail.v2"
  "html/template"
)

func parseTemplate(fileName string, data interface{}) (string, error) {
  t, err := template.ParseFiles(fileName)
  if err != nil {
    return "", err
  }

  buf := new(bytes.Buffer)
  if err = t.Execute(buf, data); err != nil {
    return "", err
  }

  return buf.String(), nil
}

func sendEmail(email string, header string, result string) error {
  m := gomail.NewMessage()
  m.SetHeader("From", configs.GetEnv("EMAIL_USERNAME"))
  m.SetHeader("To", email)
  m.SetHeader("Subject", header)
  m.SetBody("text/html", result)

  d := gomail.NewDialer("smtp.gmail.com", 587, configs.GetEnv("EMAIL_USERNAME"), configs.GetEnv("EMAIL_PASSWORD"))

  // send email
  if err := d.DialAndSend(m); err != nil {
    fmt.Println("Could not send email because: ", err)
    return err
  }

  return nil
}

func SendEmailAuth(user models.User, endpoint string, header string) {
  token, err := GenerateToken(user.ID, user.Email, user.Role.Role)
  if err != nil {
    fmt.Println("Could not generate token")
  }

  // struct for store data that used in html template
  templateData := struct {
    URL string
  }{
    URL: fmt.Sprintf("%s/api/auth/%s?token=%s", configs.GetEnv("APP_URL"), endpoint, token),
  }

  // use html templates based on endpoint
  var result string
  if endpoint == "verify-email" {
    result, err = parseTemplate("./../templates/verify-email.html", templateData)
    fmt.Println(err)
  } else {
    result, err = parseTemplate("./../templates/reset-password.html", templateData)
    fmt.Println(err)
  }

  // send email
  if err := sendEmail(user.Email, header, result); err != nil {
    fmt.Println("Could not send email because: ", err)
  }
}

func SendEmailBookingConfirmation(userBio models.UserBio, flight models.Flight, booking models.Booking, header string) {
  // struct for store data that used in html template
  templateData := struct {
    FirstName   string
    LastName    string
    Origin      string
    Destination string
    Date        string
    Time        string
    DateLimit   string
    TimeLimit   string
    PhoneNumber string
    Total       int
    BookingID   uint
  }{
    FirstName:   userBio.FirstName,
    LastName:    userBio.LastName,
    Origin:      flight.Origin,
    Destination: flight.Destination,
    Date:        flight.CreatedAt.Format("2006-01-02"),
    Time:        flight.CreatedAt.Format("15:04"),
    DateLimit:   flight.CreatedAt.AddDate(0, 0, 1).Format("2006-01-02"),
    TimeLimit:   flight.CreatedAt.AddDate(0, 0, 1).Format("15:04"),
    PhoneNumber: userBio.PhoneNumber,
    Total:       int(booking.TotalPrice),
    BookingID:   booking.ID,
  }

  // parse html template
  result, _ := parseTemplate("./../templates/booking.html", templateData)

  // send email
  if err := sendEmail(userBio.User.Email, header, result); err != nil {
    fmt.Println("Could not send email because: ", err)
  }
}

func SendEmailPaymentConfirmation(userBio models.UserBio, payment models.Payment, header string) {
  // struct for store data that used in html template
  templateData := struct {
    FirstName   string
    LastName    string
    Date        string
    Time        string
    PhoneNumber string
    Total       int
    PaymentID   uint
  }{
    FirstName:   userBio.FirstName,
    LastName:    userBio.LastName,
    Date:        payment.CreatedAt.Format("2006-01-02"),
    Time:        payment.CreatedAt.Format("15:04"),
    PhoneNumber: userBio.PhoneNumber,
    Total:       int(payment.TotalPrice),
    PaymentID:   payment.ID,
  }

  // parse html template
  result, _ := parseTemplate("./../templates/payment.html", templateData)

  // send email
  if err := sendEmail(userBio.User.Email, header, result); err != nil {
    fmt.Println("Could not send email because: ", err)
  }
}
