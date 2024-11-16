package dbutil

import "database/sql"

type ApplicationService struct {
	db *sql.DB
}

func NewApplicationService(db *sql.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

// 物理删除
// 逻辑删除: 使用 Update()
func (s *ApplicationService) Remove() error {
	return nil
}

// 逻辑删除：更新 state 列 { "deprecated": true }
func (s *ApplicationService) Update() error {
	return nil
}

func (s *ApplicationService) Get() (map[string]interface{}, error) {
	return nil, nil
}

func (s *ApplicationService) GetMany() ([]map[string]interface{}, error) {
	return nil, nil
}

func (s *ApplicationService) Create() error {
	return nil
}
