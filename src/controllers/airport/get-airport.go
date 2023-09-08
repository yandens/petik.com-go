package airport

import (
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/utils"
  "io/ioutil"
  "net/http"
)

type APIResponse struct {
  Type     string `json:"type"`
  Features []struct {
    Type     string `json:"type"`
    Geometry struct {
      Type        string    `json:"type"`
      Coordinates []float64 `json:"coordinates"`
    } `json:"geometry"`
    Properties struct {
      ID      int `json:"id"`
      Country struct {
        Code      string `json:"code"`
        Name      string `json:"name"`
        Continent string `json:"continent"`
        Wikipedia string `json:"wikipedia"`
      } `json:"country"`
      Name           string      `json:"name"`
      Source         string      `json:"source"`
      Distance       interface{} `json:"distance"`
      MatchRelevance struct {
        Code           interface{} `json:"code"`
        Country        interface{} `json:"country"`
        Levenshtein    interface{} `json:"levenshtein"`
        TSRank         interface{} `json:"ts_rank"`
        TrgmSimilarity interface{} `json:"trgm_similarity"`
        SkippedChunks  int         `json:"skipped_chunks"`
      } `json:"match_relevance"`
      MatchLevel int `json:"match_level"`
      Region     struct {
        Code      string `json:"code"`
        LocalCode string `json:"local_code"`
        Name      string `json:"name"`
        Wikipedia string `json:"wikipedia"`
      } `json:"region"`
      Elevation    int         `json:"elevation"`
      Functions    []string    `json:"functions"`
      GPSCode      string      `json:"gps_code"`
      HomeLink     string      `json:"home_link"`
      IATA         string      `json:"iata"`
      LocalCode    interface{} `json:"local_code"`
      Municipality string      `json:"municipality"`
      Type         string      `json:"type"`
      Wikipedia    string      `json:"wikipedia"`
    } `json:"properties"`
  } `json:"features"`
  Properties struct{} `json:"properties"`
}

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
  var data APIResponse
  if err := json.Unmarshal(respData, &data); err != nil {
    utils.JSONResponse(c, 400, false, "Failed to get airport", nil)
    return
  }

  // return response
  utils.JSONResponse(c, 200, true, "Airport retrieved successfully", data)
}

func GetAirportByIATA(IATA string) (interface{}, error) {
  // define url
  url := fmt.Sprintf("https://port-api.com/airport/search/%s", IATA)

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
  var data APIResponse
  if err := json.Unmarshal(respData, &data); err != nil {
    return nil, err
  }

  // return response
  return data, nil
}
