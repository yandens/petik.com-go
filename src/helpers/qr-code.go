package helpers

import (
  "encoding/base64"
  "encoding/json"
  "github.com/skip2/go-qrcode"
  "github.com/yandens/petik.com-go/src/controllers/ticket"
)

func GenerateQRCode(data ticket.TicketData) (string, error) {
  // convert data to json
  jsonData, err := json.Marshal(data)
  if err != nil {
    return "", err
  }

  // generate qr code
  qrCode, err := qrcode.Encode(string(jsonData), qrcode.Medium, 256)
  if err != nil {
    return "", err
  }

  // encode to base64
  fileBase64 := base64.StdEncoding.EncodeToString(qrCode)

  return fileBase64, nil
}
