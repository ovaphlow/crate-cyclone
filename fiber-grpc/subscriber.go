package main

import (
	"database/sql"
	"ovaphlow/cratecyclone/utility"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Subscriber struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Tags   string `json:"tags"`
	Detail string `json:"detail"`
	Time   string `json:"time"`
	State  string `json:"state"`
}

type SubscriberRepo interface {
	retrieveByID(id int64, uuid string) (*Subscriber, error)
	save(subscriber *Subscriber) error
	retrieveByUsername(userName string) (*Subscriber, error)
}

type SubscriberRepoImpl struct {
	db *sql.DB
}

func NewSubscriberRepoImpl(db *sql.DB) *SubscriberRepoImpl {
	return &SubscriberRepoImpl{db: db}
}

func (r *SubscriberRepoImpl) retrieveByID(id int64, uuid string) (*Subscriber, error) {
	subscriber := &Subscriber{}
	err := r.db.QueryRow("SELECT id, email, name, phone, tags, detail, time, state FROM crate.subscriber WHERE id = $1 AND state->>'uuid' = $2", id, uuid).Scan(&subscriber.ID, &subscriber.Email, &subscriber.Name, &subscriber.Phone, &subscriber.Tags, &subscriber.Detail, &subscriber.Time, &subscriber.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Subscriber{}, nil
		} else {
			utility.Slogger.Error(err.Error())
			return nil, err
		}
	}
	return subscriber, nil
}

func (r *SubscriberRepoImpl) save(subscriber *Subscriber) error {
	_, err := r.db.Exec("INSERT INTO subscriber (email, name, phone, tags, detail, time, state) VALUES (?, ?, ?, ?, ?, ?, ?)", subscriber.Email, subscriber.Name, subscriber.Phone, subscriber.Tags, subscriber.Detail, subscriber.Time, subscriber.State)
	if err != nil {
		return err
	}
	return nil
}

func (r *SubscriberRepoImpl) retrieveByUsername(userName string) (*Subscriber, error) {
	subscriber := &Subscriber{}
	err := r.db.QueryRow("SELECT id, email, name, phone, tags, detail, time, state FROM subscriber WHERE name = ?", userName).Scan(&subscriber.ID, &subscriber.Email, &subscriber.Name, &subscriber.Phone, &subscriber.Tags, &subscriber.Detail, &subscriber.Time, &subscriber.State)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}

type SubscriberService struct {
	repo SubscriberRepo
}

func NewSubscriberService(repo SubscriberRepo) *SubscriberService {
	return &SubscriberService{repo: repo}
}

func (s *SubscriberService) RetrieveByID(id int64, uuid string) (*Subscriber, error) {
	return s.repo.retrieveByID(id, uuid)
}

func (s *SubscriberService) Save(subscriber *Subscriber) error {
	return s.repo.save(subscriber)
}

func (s *SubscriberService) LogIn(userName string) (*Subscriber, error) {
	return s.repo.retrieveByUsername(userName)
}

func AddSubscriberEndpoints(app *fiber.App, s *SubscriberService) {
	app.Get("/crate-api/subscriber/:id", func(c *fiber.Ctx) error {
		uuid := c.Query("uuid", "")
		id, err := strconv.ParseInt(c.Params("id", "0"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
		}
		subscriber, err := s.RetrieveByID(id, uuid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
		}
		return c.Status(200).JSON(subscriber)
	})
}
