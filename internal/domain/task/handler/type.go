package handler

import "time"

type SingleTask struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type Create struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Update struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
