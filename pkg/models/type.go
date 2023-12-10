package models

import "time"

type IP struct {
	ID         int64     `json:"id"`
	Data       string    `json:"data"`
	Type1      string    `json:"type1"`
	Type2      string    `json:"type2"`
	Speed      int64     `json:"speed"`  // 链接速度
	Source     string    `json:"source"` // 代理来源
	CreateTime time.Time `json:"create_time"`
}
