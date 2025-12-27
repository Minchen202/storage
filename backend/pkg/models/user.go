package models

import (
	"time"
)

type User struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"` // Do not expose this in JSON responses
	CreatedAt          time.Time `json:"created_at"`
	StorageContributed int64     `json:"storage_contributed"`
	StorageQuota       int64     `json:"storage_quota"`
	TotalUptimeSeconds int64     `json:"total_uptime_seconds"`
	LastOnline         time.Time `json:"last_online"`
}
