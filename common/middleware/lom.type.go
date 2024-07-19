package middleware

import "time"

type LOMKeyWithData struct {
	Data       string    `json:"data"`
	ExpireDate time.Time `json:"expire_date"`
}
