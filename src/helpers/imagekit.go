package helpers

import (
  "encoding/base64"
  "github.com/codedius/imagekit-go"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
  "io"
  "mime/multipart"
)

func UploadToImagekit(c *gin.Context, file *multipart.FileHeader, fileName string) (string, error) {
  // construct imagekit
  opts := imagekit.Options{
    PublicKey:  configs.GetEnv("IMAGEKIT_PUBLIC_KEY"),
    PrivateKey: configs.GetEnv("IMAGEKIT_PRIVATE_KEY"),
  }

  ik, err := imagekit.NewClient(&opts)
  if err != nil {
    return "", err
  }

  // open file
  f, err := file.Open()
  if err != nil {
    return "", err
  }
  defer f.Close()

  // read file
  fileBytes, err := io.ReadAll(f)
  if err != nil {
    return "", err
  }

  // encode file into base64
  fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)

  // upload avatar
  ur := imagekit.UploadRequest{
    File:     fileBase64,
    FileName: fileName,
    Folder:   "/avatars",
  }

  upr, err := ik.Upload.ServerUpload(c, &ur)
  if err != nil {
    return "", err
  }

  return upr.URL, nil
}
