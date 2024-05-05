package subscriber

import "time"

type Subscriber struct {
	ID     int64     `json:"id"`
	Time   time.Time `json:"time"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Phone  string    `json:"phone"`
	Tags   string    `json:"tags"`
	Detail string    `json:"detail"`
	State  string    `json:"state"`
}
