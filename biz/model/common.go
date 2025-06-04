package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type CommonBase struct {
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
	IsDelete  IsDelete   `gorm:"column:is_delete"`
}

type IsDelete int

const (
	DeleteTrue  IsDelete = 1
	DeleteFalse IsDelete = 0
)

type StringArray []string

func (m *StringArray) Scan(val interface{}) error {
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal value:", val))
	}
	var result []string
	err := json.Unmarshal(bytes, &result)
	*m = result
	return err
}

func (m StringArray) Value() (driver.Value, error) {
	str, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return str, nil
}

type StringMap map[string]interface{}

func (m *StringMap) Scan(val interface{}) error {
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal value:", val))
	}
	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	*m = result
	return err
}

func (m StringMap) Value() (driver.Value, error) {
	str, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return str, nil
}
