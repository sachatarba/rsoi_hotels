package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ReservationStatus int

const (
	Paid ReservationStatus = iota
	Canceled
)

func (s ReservationStatus) String() string {
	switch s {
	case Paid:
		return "PAID"
	case Canceled:
		return "CANCELED"
	default:
		return "UNKNOWN"
	}
}

func NewReservationStatus(status string) ReservationStatus {
	switch status {
	case "PAID":
		return Paid
	case "CANCELED":
		return Canceled
	default:
		return -1
	}
}

func (s ReservationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s ReservationStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *ReservationStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var strVal string
	switch v := value.(type) {
	case string:
		strVal = v
	case []byte:
		strVal = string(v)
	default:
		return fmt.Errorf("failed to scan ReservationStatus: expected string or []byte, got %T", value)
	}

	*s = NewReservationStatus(strVal)
	return nil
}
