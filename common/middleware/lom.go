package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/Sincerelyzl/larb-on-me/common/utils"
)

var (
	// LOMSecretKey        = "okdisodkfi2kdisk"
	// LOMExpireTime       = time.Minute * 10
	// LOMCookieAuthPrefix = "LOM-AUTH"
	// LOMUserPrefix       = "LOM-USER"
	LOMSecretKey        = utils.EnvString("LOM_SECRET_KEY", "okdisodkfi2kdisk")
	LOMExpireTime       = utils.EnvDuration("LOM_EXPIRE_TIME", "10m")
	LOMCookieAuthPrefix = utils.EnvString("LOM_COOKIE_AUTH_PREFIX", "LOM-AUTH")
	LOMUserPrefix       = utils.EnvString("LOM_USER_PREFIX", "LOM-USER")
)

func encryptString(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(LOMSecretKey))
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

	block, err := aes.NewCipher([]byte(LOMSecretKey))
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
		ExpireDate: time.Now().Add(LOMExpireTime),
	}
	LogGlobal.Log.Info("info", "expire time", LOMExpireTime)

	bytesKey, err := json.Marshal(key)
	if err != nil {
		return "", err
	}

	jsonString := string(bytesKey)

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
