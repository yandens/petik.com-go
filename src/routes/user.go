package routes

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/controllers/bio"
  "github.com/yandens/petik.com-go/src/controllers/booking"
  "github.com/yandens/petik.com-go/src/controllers/notification"
  "github.com/yandens/petik.com-go/src/controllers/payment"
  "github.com/yandens/petik.com-go/src/middlewares"
)

func UserRoutes(router *gin.RouterGroup) {
  userRoute := router.Group("/user")
  userRoute.Use(middlewares.Authorized("user"))

  // bio routes
  userRoute.POST("/create-bio", bio.CreateBio)
  userRoute.PUT("/update-bio", bio.UpdateBio)
  userRoute.GET("/get-bio", bio.ReadBio)
  userRoute.POST("/upload-avatar", bio.UploadAvatar)

  // booking routes
  userRoute.POST("/create-booking", booking.CreateBooking)
  userRoute.PUT("/cancel-booking", booking.CancelBooking)
  userRoute.GET("/my-booking", booking.GetBookings)

  // payment routes
  userRoute.POST("/create-payment", payment.CreatePayment)

  // notification routes
  userRoute.GET("/get-notifications", notification.GetNotifications)
}
