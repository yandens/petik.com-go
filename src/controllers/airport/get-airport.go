package airport

import (
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
  "net/http"
)

func GetAirport(c *gin.Context) {
  // get params
  search := c.Param("search")
  if search == "" {
    utils.JSONResponse(c, 400, false, "Search is required", nil)
  }

  // define url
  url := "https://port-api.com/airport/search/" + search

  // call third party api to get airport
  resp, err := http.Get(url)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Failed to get airport", nil)
    return
  }
  defer resp.Body.Close()

  // return response
  utils.JSONResponse(c, 200, true, "Airport retrieved successfully", resp.Body)
}
