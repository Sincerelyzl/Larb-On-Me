package models

import "time"

type LOMKeyWithData struct {
	Data       string    `json:"data"`
	ExpireDate time.Time `json:"expire_date"`
}

type LOMUser struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
