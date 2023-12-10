package helpers

import (
  "github.com/codedius/imagekit-go"
  "github.com/gin-gonic/gin"
  "github.com/yandens/petik.com-go/src/configs"
)

func UploadToImagekit(c *gin.Context, file string, fileName string, folder string) (string, error) {
  // construct imagekit
  opts := imagekit.Options{
    PublicKey:  configs.GetEnv("IMAGEKIT_PUBLIC_KEY"),
    PrivateKey: configs.GetEnv("IMAGEKIT_PRIVATE_KEY"),
  }

  ik, err := imagekit.NewClient(&opts)
  if err != nil {
    return "", err
  }

  // upload avatar
  ur := imagekit.UploadRequest{
    File:     file,
    FileName: fileName,
    Folder:   folder,
  }

  upr, err := ik.Upload.ServerUpload(c, &ur)
  if err != nil {
    return "", err
  }

  return upr.URL, nil
}
