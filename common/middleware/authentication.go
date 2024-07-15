package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	superSecretKey = "okdisodkfi2kdisk"
	expireTime     = time.Second * 3
)

type LOMKeyWithData struct {
	Data       string    `json:"data"`
	ExpireDate time.Time `json:"expire_date"`
}

type LOMKeyWithOutExpire struct {
	Data interface{} `json:"data"`
}

func encryptString(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(superSecretKey))
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	encryptedData := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func decryptString(encryptedText string) (*LOMKeyWithData, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(superSecretKey))
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("LOM key is invalid")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("LOM Authentication failed")
	}

	LOMData := LOMKeyWithData{}
	errUnmarshal := json.Unmarshal(plaintext, &LOMData)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &LOMData, nil
}

func GenerateLOMKeys(data any) (string, error) {

	bytesData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	jsonData := string(bytesData)
	key := LOMKeyWithData{
		Data:       jsonData,
		ExpireDate: time.Now().Add(expireTime),
	}

	bytesKey, err := json.Marshal(key)
	if err != nil {
		return "", err
	}

	jsonString := string(bytesKey)

	fmt.Println("before-encrypt:" + jsonString)
	encrypted, err := encryptString(jsonString)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func ClaimsLOM(encrypted string, toParse any) error {

	decrypted, err := decryptString(encrypted)
	if err != nil {
		return err
	}

	if decrypted.ExpireDate.Before(time.Now()) {
		return errors.New("token has expired")
	}

	errUnmarshal := json.Unmarshal([]byte(decrypted.Data), toParse)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	return nil
}
