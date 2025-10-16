package model

import "time"

type Duty struct {
	Date  time.Time
	Hours int
	Who   string
}
