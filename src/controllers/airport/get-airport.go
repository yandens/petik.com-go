package airport

import (
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
  "io/ioutil"
  "net/http"
)

func GetAirportBySearch(c *gin.Context) {
  // get params
  search := c.Param("search")
  if search == "" {
    utils.JSONResponse(c, 400, false, "Search is required", nil)
  }

  // define url
  url := fmt.Sprintf("https://port-api.com/airport/search/%s", search)

  // call third party api to get airport
  resp, err := http.Get(url)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Failed to get airport", nil)
    return
  }
  defer resp.Body.Close()

  // read response
  respData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    utils.JSONResponse(c, 400, false, "Failed to get airport", nil)
    return
  }

  // parse response to readable json
  var data map[string]interface{}
  if err := json.Unmarshal(respData, &data); err != nil {
    utils.JSONResponse(c, 400, false, "Failed to get airport", nil)
    return
  }

  // return response
  utils.JSONResponse(c, 200, true, "Airport retrieved successfully", data)
}

func GetAirportByIATA(IATA string) (interface{}, error) {
  // define url
  url := fmt.Sprintf("https://port-api.com/airport/iata/%s", IATA)

  // call third party api to get airport
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }

  // read response
  respData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  // parse response to readable json
  var data map[string]interface{}
  if err := json.Unmarshal(respData, &data); err != nil {
    return nil, err
  }

  // return response
  return data, nil
}
