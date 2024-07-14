package file

import "time"

type File struct {
	ID    string    `json:"id"`
	State string    `json:"state"`
	Time  time.Time `json:"time"`
}
