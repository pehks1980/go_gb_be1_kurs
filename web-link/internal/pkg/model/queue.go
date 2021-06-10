package model

import "time"

type Data struct {
	Data []DataEl `json:"data"`
}

// Data элемент строки файла json
type DataEl struct {
	UID      string    `json:"uid"`
	URL      string    `json:"url"`
	Shorturl string    `json:"shorturl"`
	Datetime time.Time `json:"datetime"`
	Active   int       `json:"active"`
	Redirs   int       `json:"redirs"`
}
