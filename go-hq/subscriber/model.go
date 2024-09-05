package subscriber

import "time"

type Subscriber struct {
	ID          string    `json:"id"`
	State       string    `json:"state"`
	Time        time.Time `json:"time"`
	RelationID  string    `json:"relationId"`
	ReferenceID string    `json:"referenceId"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Tags        string    `json:"tags"`
	Detail      string    `json:"detail"`
}
