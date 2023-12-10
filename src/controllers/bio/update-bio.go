package bio

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/models"
)

type UpdateBioInput struct {
  FirstName   string `json:"firstName"`
  LastName    string `json:"lastName"`
  PhoneNumber string `json:"phoneNumber" binding:"min=12,max=13"`
  Address     string `json:"address"`
  City        string `json:"city"`
  Province    string `json:"province"`
  Avatar      string `json:"avatar"`
}

func UpdateBio(c *gin.Context) {
  // connect to database
  db, err := configs.ConnectToDB()
  if err != nil {
    helpers.JSONResponse(c, 500, false, "Could not connect to the database", nil)
    return
  }

  // get user id
  id, _ := c.Get("id")
  if id == "" {
    helpers.JSONResponse(c, 401, false, "Unauthorized", nil)
    return
  }

  // get user input
  var input UpdateBioInput
  if err := c.ShouldBindJSON(&input); err != nil {
    helpers.JSONResponse(c, 400, false, "Input must be JSON", nil)
    return
  }

  // update user bio
  if err := db.Model(&models.UserBio{}).Where("user_id = ?", id).Updates(input).Error; err != nil {
    helpers.JSONResponse(c, 400, false, "Failed to update user bio", nil)
    return
  }

  helpers.JSONResponse(c, 200, true, "User bio updated successfully", nil)
}
