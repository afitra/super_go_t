package model

import "time"

type Chache_model struct {
	Key     string
	Data    []byte
	Expired time.Duration
}
