package subscriber

import "database/sql"

type SubscriberRepo interface {
	getByParameter(p string) (*Subscriber, error)
}

type SubscriberRepoImpl struct {
	db *sql.DB
}

func NewSubscriberRepoImpl(db *sql.DB) *SubscriberRepoImpl {
	return &SubscriberRepoImpl{db: db}
}

func (r *SubscriberRepoImpl) getByParameter(p string) (*Subscriber, error) {
	q := `
	select id, state, time, relation_id, reference_id, email, name, phone, tags, detail
	from crate.subscriber
	where email = $1 or name = $1 or phone = $1
	limit 1
	`
	var s Subscriber
	err := r.db.QueryRow(q, p).Scan(
		&s.ID,
		&s.State,
		&s.Time,
		&s.RelationID,
		&s.ReferenceID,
		&s.Email,
		&s.Name,
		&s.Phone,
		&s.Tags,
		&s.Detail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
