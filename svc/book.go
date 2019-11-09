package svc

import (
	"errors"
	"time"
)

var (
	ErrCheckoutInvalid = errors.New("checkout time must be at least one day after check in")
)

type Booker interface {
	Book(employee EmployeeID, hotel HotelID, room RoomType, checkIn time.Time, checkOut time.Time) error
}

func NewBooker() Booker {
	return &bookingService{}
}

type bookingService struct {
}

func (s *bookingService) Book(employee EmployeeID, hotel HotelID, room RoomType, checkIn time.Time, checkOut time.Time) error {
	stayDuration := checkOut.Sub(checkIn)
	if stayDuration <= 0 || (checkOut.Day() == checkIn.Day() && stayDuration < 24*time.Hour) {
		return ErrCheckoutInvalid
	}
	return nil
}
