package utils

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "github.com/yandens/petik.com-go/src/configs"
  "io"
)

func Encrypt(data []byte) ([]byte, error) {
  // create new cipher block based on the key
  block, err := aes.NewCipher([]byte(configs.GetEnv("ENCRYPTION_KEY")))
  if err != nil {
    return nil, err
  }

  // create new gcm based on the cipher block
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    return nil, err
  }

  // create nonce with the same size as the gcm nonce size
  nonce := make([]byte, gcm.NonceSize())
  if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
    return nil, err
  }

  // encrypt the data with the nonce
  ciphertext := gcm.Seal(nonce, nonce, data, nil)

  return ciphertext, nil
}

func Decrypt(data []byte) ([]byte, error) {
  // create new cipher block based on the key
  block, err := aes.NewCipher([]byte(configs.GetEnv("ENCRYPTION_KEY")))
  if err != nil {
    return nil, err
  }

  // create new gcm based on the cipher block
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    return nil, err
  }

  // get the nonce size from the gcm
  nonceSize := gcm.NonceSize()
  // get the nonce from the data
  nonce, ciphertext := data[:nonceSize], data[nonceSize:]

  // decrypt the data with the nonce
  plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
  if err != nil {
    return nil, err
  }

  return plaintext, nil
}
