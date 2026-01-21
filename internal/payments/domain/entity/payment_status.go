package entity

type PaymentStatus string

const (
	Paid     PaymentStatus = "PAID"
	Canceled PaymentStatus = "CANCELED"
)

func (s PaymentStatus) String() string {
	return string(s)
}
