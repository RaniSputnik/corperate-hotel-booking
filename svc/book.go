package svc

import (
	"errors"
	"time"
)

var (
	ErrCheckoutInvalid = errors.New("checkout time must be at least one day after check in")
	ErrInvalidHotel    = errors.New("the given hotel id is not valid")
	ErrRoomNotOfferred = errors.New("the hotel selected does not offer that room type")
	ErrNotAvailable    = errors.New("not available")
)

type booking struct {
	Hotel      HotelID
	RoomNumber string
	CheckIn    time.Time
	CheckOut   time.Time
}

func NewBooker(hotels Hotels) Booker {
	return &bookingService{
		hotels:   hotels,
		bookings: []*booking{},
	}
}

type bookingService struct {
	hotels   Hotels
	bookings []*booking
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

	availableRooms := hotel.Options(room)
	existingBookings := getBookingsForPeriod(hotelID, s.bookings, checkIn, checkOut)
	if len(existingBookings) > len(availableRooms)-1 {
		return ErrNotAvailable
	}

	// TODO: Do we need to modify check in and out times so that
	// they are for the correct hour for the given hotel
	thisBooking := &booking{
		Hotel:      hotelID,
		RoomNumber: "TODO",
		CheckIn:    checkIn,
		CheckOut:   checkOut,
	}
	s.bookings = append(s.bookings, thisBooking)

	return nil
}

func getBookingsForPeriod(hotelID HotelID, bookings []*booking, checkIn, checkOut time.Time) []*booking {
	results := []*booking{}
	for _, booking := range bookings {
		if booking.Hotel != hotelID {
			continue
		}
		if booking.CheckIn.After(checkOut) || booking.CheckOut.Before(checkIn) {
			continue
		}
		results = append(results, booking)
	}
	return results
}
