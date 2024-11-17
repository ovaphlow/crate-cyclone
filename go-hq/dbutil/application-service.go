package dbutil

import (
	"encoding/json"
	"ovaphlow/crate/hq/utility"
	"time"
)

type ApplicationService interface {
	Create(st string, d map[string]interface{}) error
	Get(st string, f [][]string, l string) (map[string]interface{}, error)
	Update(st string, d map[string]interface{}, w string, deprecated bool) error
	Remove(st string, w string) error
}

type ApplicationServiceImpl struct {
	repo *SharedRepoImpl
}

func NewApplicationService(repo *SharedRepoImpl) *ApplicationServiceImpl {
	return &ApplicationServiceImpl{repo: repo}
}

func (s *ApplicationServiceImpl) Create(st string, d map[string]interface{}) error {
	// id
	id, err := utility.GenerateKsuid()
	if err != nil {
		return err
	}
	d["id"] = id

	time_string := time.Now().Format("2006-01-02 15:04:05-0700")

	// time
	d["time"] = time_string

	// state
	state := map[string]interface{}{
		"created_at": time_string,
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return err
	}
	d["state"] = string(stateJson)

	return s.repo.Create(st, d)
}

func (s *ApplicationServiceImpl) Get(st string, f [][]string, l string) (map[string]interface{}, error) {
	return nil, nil
}

// 修改
// 逻辑删除：更新 state 列 { "deprecated": true }
func (s *ApplicationServiceImpl) Update(st string, d map[string]interface{}, w string, deprecated bool) error {
	// state
	state := map[string]interface{}{
		"updated_at": time.Now().Format("2006-01-02 15:04:05-0700"),
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return err
	}
	d["state"] = string(stateJson)

	return s.repo.Update(st, d, w)
}

// 物理删除
// 逻辑删除: 使用 Update()
func (s *ApplicationServiceImpl) Remove(st string, w string) error {
	return nil
}
