package model

import "time"

type Message struct {
	Msg       string
	CreatedAt time.Time
}

type UserBlock struct {
	BlockUser string
}

type BlockList struct {
	BlockList []string
}
