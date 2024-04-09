package domain

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type IdPagination struct {
	AfterID string `json:"after_id"`
}

type IntTime time.Time

func (it *IntTime) Scan(src interface{}) error {
	t, ok := src.(time.Time)
	if !ok {
		return errors.New("cannot assert time from db to type time.Time")
	}
	*it = IntTime(t)

	return nil
}

func (it *IntTime) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*it = IntTime(time.UnixMilli(n))

	return nil
}

func (it IntTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(it).UnixMilli())
}
