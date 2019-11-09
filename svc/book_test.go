package svc_test

import (
	"testing"
	"time"

	"github.com/RaniSputnik/corperate-hotel-booking/svc"
	"github.com/stretchr/testify/assert"
)

const anyEmployee = svc.EmployeeID("some-employee")
const anyHotel = svc.HotelID("some-hotel")
const anyRoom = svc.RoomSingle

var anyCheckIn, _ = time.Parse(shortForm, "2019-Apr-22")
var anyCheckOut, _ = time.Parse(shortForm, "2019-Apr-30")

const shortForm = "2006-Jan-02"
const longForm = "2006-Jan-02 at 3:04pm"

func TestBookDates(t *testing.T) {
	mockHotels := &HotelsMock{}
	mockHotels.Func.GetHotel.Returns.Hotel = &svc.Hotel{
		ID:   "some-hotel",
		Name: "Some Hotel",
		Rooms: []svc.Room{
			svc.Room{Type: anyRoom, Number: "04"},
		},
	}
	booker := svc.NewBooker(mockHotels)

	t.Run("fails if check-out date matches check-in", func(t *testing.T) {
		checkIn, _ := time.Parse(shortForm, "2019-Feb-03")
		checkOut := checkIn

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("fails if the check-out date is later the same day", func(t *testing.T) {
		checkIn, _ := time.Parse(longForm, "2019-Feb-03 at 9:00am")
		checkOut, _ := time.Parse(longForm, "2019-Feb-03 at 2:00pm")

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("fails if check-out date is before check-in", func(t *testing.T) {
		checkIn, _ := time.Parse(shortForm, "2019-Feb-03")
		checkOut, _ := time.Parse(shortForm, "2019-Feb-02")

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("succeeds if the check-out date is one day after check-in", func(t *testing.T) {
		checkIn, _ := time.Parse(shortForm, "2019-Feb-03")
		checkOut, _ := time.Parse(shortForm, "2019-Feb-04")

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.NoError(t, err)
	})

	t.Run("succeeds if the check-out date is one month after checkin", func(t *testing.T) {
		checkIn, _ := time.Parse(shortForm, "2019-Feb-03")
		checkOut, _ := time.Parse(shortForm, "2019-Mar-03")

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.NoError(t, err)
	})

	t.Run("ignores the hours minutes component of the times", func(t *testing.T) {
		checkIn, _ := time.Parse(longForm, "2019-Feb-03 at 11:59pm")
		checkOut, _ := time.Parse(longForm, "2019-Feb-04 at 12:01am")

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.NoError(t, err)
	})
}

func TestBookRoom(t *testing.T) {
	t.Run("fails if the hotel does not exist", func(t *testing.T) {
		noHotels := &HotelsMock{}
		noHotels.Func.GetHotel.Returns.Hotel = nil
		noHotels.Func.GetHotel.Returns.Err = nil
		booker := svc.NewBooker(noHotels)

		const invalidHotel = svc.HotelID("The Infinity Hotel")
		err := booker.Book(anyEmployee, invalidHotel, anyRoom, anyCheckIn, anyCheckOut)
		assert.Equal(t, svc.ErrInvalidHotel, err)
	})

	t.Run("fails if the hotel does not offer the given room type", func(t *testing.T) {
		theHotel := svc.HotelID("the-hotel")
		singleRooms := &HotelsMock{}
		singleRooms.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   theHotel,
			Name: "The Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "04"},
			},
		}
		booker := svc.NewBooker(singleRooms)

		const invalidHotel = svc.HotelID("The Infinity Hotel")
		err := booker.Book(anyEmployee, theHotel, svc.RoomDouble, anyCheckIn, anyCheckOut)
		assert.Equal(t, svc.ErrRoomNotOfferred, err)
	})

	t.Run("fails if the room type is invalid", func(t *testing.T) {
		singleRooms := &HotelsMock{}
		singleRooms.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   anyHotel,
			Name: "Any Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "04"},
			},
		}
		booker := svc.NewBooker(singleRooms)

		err := booker.Book(anyEmployee, anyHotel, svc.RoomType("something-bogus"), anyCheckIn, anyCheckOut)
		assert.Equal(t, svc.ErrRoomNotOfferred, err)
	})
}
