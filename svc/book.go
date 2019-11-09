package svc

import (
	"errors"
	"time"
)

var (
	ErrCheckoutInvalid = errors.New("checkout time must be at least one day after check in")
	ErrInvalidHotel    = errors.New("the given hotel id is not valid")
	ErrRoomNotOfferred = errors.New("the hotel selected does not offer that room type")
)

func NewBooker(hotels Hotels) Booker {
	return &bookingService{
		hotels: hotels,
	}
}

type bookingService struct {
	hotels Hotels
}

func (s *bookingService) Book(employee EmployeeID, hotelID HotelID, room RoomType, checkIn time.Time, checkOut time.Time) error {
	stayDuration := checkOut.Sub(checkIn)
	if stayDuration <= 0 || (checkOut.Day() == checkIn.Day() && stayDuration < 24*time.Hour) {
		return ErrCheckoutInvalid
	}

	hotel, err := s.hotels.GetHotel(hotelID)
	if err != nil {
		return err
	} else if hotel == nil {
		return ErrInvalidHotel
	}

	if !hotel.Offers(room) {
		return ErrRoomNotOfferred
	}

	return nil
}
