package entity

import (
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain"
)

type Loyalty struct {
	Id               int
	Username         string
	reservationCount int
	status           LoyaltyStatus
	discount         int
}

func NewLoyalty(id int, username string, reservationCount int,
	status LoyaltyStatus, discount int) (*Loyalty, error) {
	realStatus := getStatus(reservationCount)
	if realStatus != status {
		return nil, domain.ErrLoyaltyStatusNotRespondReservations
	}

	realDiscount := getDiscount(status)
	if realDiscount != discount {
		return nil, domain.ErrDiscountNotRespondLoyaltyStatus
	}

	return &Loyalty{
		Id:               id,
		Username:         username,
		reservationCount: reservationCount,
		status:           status,
		discount:         discount,
	}, nil
}

func (l *Loyalty) ReservationCount() int {
	return l.reservationCount
}

func (l *Loyalty) Status() LoyaltyStatus {
	return l.status
}

func (l *Loyalty) Discount() int {
	return l.discount
}

func (l *Loyalty) AddReservations(reservationsCount int) {
	l.reservationCount += reservationsCount
	newStatus := getStatus(l.reservationCount)
	newDiscount := getDiscount(newStatus)

	l.status = newStatus
	l.discount = newDiscount
}
