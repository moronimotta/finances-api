package meta

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Meta map[string]string

func New() Meta {
	return make(map[string]string)
}

func (m *Meta) Add(key, value string) {
	(*m)[key] = value
}

func (m *Meta) Delete(key string) {
	if m == nil {
		return
	}

	delete((*m), key)
}

func (m *Meta) Get(key string) string {
	return (*m)[key]
}

func (m *Meta) Merge(meta Meta) Meta {
	if m == nil {
		return meta
	}

	if len(meta) == 0 {
		return *m
	}

	for k, v := range meta {
		(*m)[k] = v
	}

	return *m
}

func (m Meta) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan scans value into Jsonb and implements sql.Scanner interface
func (m *Meta) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}
