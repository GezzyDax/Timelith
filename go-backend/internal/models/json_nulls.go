package models

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

type NullString struct {
	sql.NullString
}

func NewNullString(value string) NullString {
	return NullString{sql.NullString{String: value, Valid: true}}
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if stringNull(data) {
		ns.String = ""
		ns.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &ns.String); err != nil {
		return err
	}

	ns.Valid = true
	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

func NewNullInt64(value int64) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: value, Valid: true}}
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	if stringNull(data) {
		ni.Int64 = 0
		ni.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &ni.Int64); err != nil {
		return err
	}

	ni.Valid = true
	return nil
}

type NullTime struct {
	sql.NullTime
}

func NewNullTime(value time.Time) NullTime {
	return NullTime{sql.NullTime{Time: value, Valid: true}}
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time.Format(time.RFC3339))
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if stringNull(data) {
		nt.Time = time.Time{}
		nt.Valid = false
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}

	nt.Time = t
	nt.Valid = true
	return nil
}

func stringNull(data []byte) bool {
	return strings.EqualFold(string(data), "null")
}
