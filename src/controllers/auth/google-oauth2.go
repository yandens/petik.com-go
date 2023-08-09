package auth

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/models"
  "github.com/yandens/petik.com-go/src/utils"
  "golang.org/x/oauth2"
  oauthapi "google.golang.org/api/oauth2/v2"
)

// make variable for oauth2 config
var (
  googleOauth2Config = &oauth2.Config{
    RedirectURL:  configs.GetEnv("GOOGLE_REDIRECT_URL"),
    ClientID:     configs.GetEnv("GOOGLE_CLIENT_ID"),
    ClientSecret: configs.GetEnv("GOOGLE_CLIENT_SECRET"),
    Scopes: []string{
      "https://www.googleapis.com/auth/userinfo.email",
      "https://www.googleapis.com/auth/userinfo.profile",
    },
    Endpoint: oauth2.Endpoint{
      AuthURL:  "https://accounts.google.com/o/oauth2/auth",
      TokenURL: "https://accounts.google.com/o/oauth2/token",
    },
  }
)

func GoogleOauth2(c *gin.Context) {
  // get url from google oauth2 config
  url := googleOauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)

  // redirect to url
  c.Redirect(302, url)
}

func GoogleOauth2Callback(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    utils.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get code from query string
  code := c.Query("code")
  if code == "" {
    utils.JSONResponse(c, 400, false, "Code is required", nil)
    return
  }

  // exchange code to token
  token, err := googleOauth2Config.Exchange(c, code)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Invalid code", nil)
    return
  }

  // get user info
  oauth2Service, _ := oauthapi.New(googleOauth2Config.Client(c, token))
  userInfo, err := oauth2Service.Userinfo.Get().Do()
  if err != nil {
    utils.JSONResponse(c, 400, false, "Invalid token", nil)
    return
  }

  // get role
  var role models.Role
  if err := db.Model(&models.Role{}).Where("role = ?", "user").First(&role).Error; err != nil {
    utils.JSONResponse(c, 400, false, "Role not found", nil)
    return
  }

  // check if user is already registered
  var user models.User
  if err := db.Model(&models.User{}).Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
    utils.JSONResponse(c, 400, false, "User already registered but using basic way (not using oauth2), please login with that way", nil)
  }

  // create user
  newUser := models.User{
    Email:       userInfo.Email,
    Password:    "",
    RoleID:      role.ID,
    AccountType: "google",
    IsVerified:  true,
  }
  if err := db.Create(&newUser).Error; err != nil {
    utils.JSONResponse(c, 500, false, "Could not create user", nil)
    return
  }

  // generate token
  authToken, err := utils.GenerateToken(newUser.ID, newUser.Email, role.Role)

  // return response
  utils.JSONResponse(c, 200, true, "Success", gin.H{
    "email": newUser.Email,
    "token": authToken,
  })
}
