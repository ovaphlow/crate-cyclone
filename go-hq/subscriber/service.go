package subscriber

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"ovaphlow/crate/hq/utility"
	"time"

	"github.com/golang-jwt/jwt"
)

type SubscriberService struct {
	repo          SubscriberRepo
	schemaService *utility.SchemaService
}

func NewSubscriberService(repo SubscriberRepo, schemaService *utility.SchemaService) *SubscriberService {
	return &SubscriberService{repo: repo, schemaService: schemaService}
}

func (s *SubscriberService) SignUp(subscriber *Subscriber) error {
	existingSubscriber, err := s.repo.getByParameter(subscriber.Email)
	if err != nil {
		return err
	}
	if existingSubscriber != nil {
		return errors.New("用户已存在")
	}
	bytes := make([]byte, 8)
	_, err = rand.Read(bytes)
	if err != nil {
		return err
	}
	key := []byte(hex.EncodeToString(bytes))
	r := hmac.New(sha256.New, key)
	r.Write([]byte(subscriber.Detail))
	sha := hex.EncodeToString(r.Sum(nil))
	detail := fmt.Sprintf(`{"salt": "%s", "sha": "%s"}`, key, sha)
	data := map[string]interface{}{
		"email":        subscriber.Email,
		"name":         subscriber.Name,
		"phone":        "",
		"tags":         "[]",
		"detail":       detail,
		"relation_id":  "",
		"reference_id": "",
	}
	err = s.schemaService.Save("crate", "subscriber", data)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriberService) LogIn(username, password, secret string) (string, error) {
	subscriber, err := s.repo.getByParameter(username)
	if err != nil {
		return "", err
	}
	if subscriber == nil {
		return "", errors.New("用户不存在")
	}
	var detail map[string]interface{}
	if err := json.Unmarshal([]byte(subscriber.Detail), &detail); err != nil {
		return "", err
	}
	salt, ok := detail["salt"].(string)
	if !ok {
		return "", err
	}
	key := []byte(salt)
	r := hmac.New(sha256.New, key)
	r.Write([]byte(password))
	sha := hex.EncodeToString(r.Sum(nil))
	if sha != detail["sha"] {
		return "", errors.New("用户名或密码错误")
	}
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Issuer:    "crate-hq",
		Subject:   subscriber.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
