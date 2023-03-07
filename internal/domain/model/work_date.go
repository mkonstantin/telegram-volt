package model

import "time"

const (
	StatusWait   = 0
	StatusAccept = 1
	StatusDone   = 2
)

type WorkDate struct {
	ID       int64
	Status   int64
	WorkDate time.Time
}
