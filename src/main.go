package main

import (
  "fmt"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  socketio "github.com/googollee/go-socket.io"
  "github.com/yandens/petik.com-go/src/configs"
  "github.com/yandens/petik.com-go/src/helpers"
  "github.com/yandens/petik.com-go/src/routes"
)

func main() {
  //db.MigrateAll()
  //db.RollbackAll()
  //db.SeedAll()
  //flight.CreateFlightSeeder()
  router := gin.Default()

  router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    MaxAge:       600,
  }))

  // routes
  // default route
  router.GET("/", func(c *gin.Context) {
    helpers.JSONResponse(c, 200, true, "Server Running Well", nil)
  })

  // api group route
  api := router.Group("/api")

  // auth routes
  routes.AuthRoutes(api)

  // user routes
  routes.UserRoutes(api)

  // admin routes
  routes.AdminRoutes(api)

  // airport routes
  routes.AirportRoutes(api)

  // flight routes
  routes.FlightRoutes(api)

  // socket io
  server := socketio.NewServer(nil)

  // connect to socket io
  server.OnConnect("/", func(s socketio.Conn) error {
    s.SetContext("")
    fmt.Println("connected:", s.ID())
    return nil
  })

  // load notifications event
  server.OnEvent("/", "load-notifications", func(s socketio.Conn, userID int) {
    // get notifications
    notifications := helpers.LoadNotifications(userID)

    // emit notification
    eventName := fmt.Sprintf("notification:%d", userID)
    s.Emit(eventName, notifications)
  })

  // read notifications event
  server.OnEvent("/", "read-notifications", func(s socketio.Conn, userID int) {
    // update notifications
    helpers.UpdateNotifications(userID)

    // get notifications
    notifications := helpers.LoadNotifications(userID)

    // emit notification
    eventName := fmt.Sprintf("notification:%d", userID)
    s.Emit(eventName, notifications)
  })

  router.Run(configs.GetEnv("HOST") + ":" + configs.GetEnv("PORT"))
}
