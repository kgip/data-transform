package po

import (
	"time"
)

type Base struct {
	ID uint
	//ID string
	CreatedAt time.Time             `gorm:"not null"`
}
