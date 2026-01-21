package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type LoyaltyStatus int

const (
	Bronze LoyaltyStatus = iota
	Silver
	Gold
)

func (l LoyaltyStatus) String() string {
	switch l {
	case Bronze:
		return "BRONZE"
	case Silver:
		return "SILVER"
	case Gold:
		return "GOLD"
	}
	return "UNKNOWN"
}

func NewLoyaltyStatus(loyaltyStatus string) LoyaltyStatus {
	switch loyaltyStatus {
	case "BRONZE":
		return Bronze
	case "SILVER":
		return Silver
	case "GOLD":
		return Gold
	}
	return Bronze
}

func (l LoyaltyStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l LoyaltyStatus) Value() (driver.Value, error) {
	return l.String(), nil
}

func (l *LoyaltyStatus) Scan(value interface{}) error {
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
		return fmt.Errorf("failed to scan LoyaltyStatus: expected string or []byte, got %T", value)
	}
	*l = NewLoyaltyStatus(strVal)
	return nil
}

func getStatus(reservationsCount int) LoyaltyStatus {
	if reservationsCount >= 20 {
		return Gold
	}
	if reservationsCount >= 10 {
		return Silver
	}
	return Bronze
}
