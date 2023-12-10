package models

import "time"

func NewIp() *IP {
	return &IP{
		Speed:      -1,
		CreateTime: time.Now(),
	}
}
