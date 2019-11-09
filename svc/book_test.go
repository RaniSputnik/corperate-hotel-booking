package svc_test

import (
	"strings"
	"testing"
	"time"

	"github.com/RaniSputnik/corperate-hotel-booking/svc"
	"github.com/stretchr/testify/assert"
)

const anyEmployee = svc.EmployeeID("some-employee")
const anyHotel = svc.HotelID("some-hotel")
const anyRoom = svc.RoomSingle

var anyCheckIn, anyCheckOut = times("2019-Apr-22", "2019-Apr-30")

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

	t.Run("fails if check-out date matches check-in", func(t *testing.T) {
		checkIn, checkOut := times("2019-Feb-03", "2019-Feb-03")
		err := svc.NewBooker(mockHotels).Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("fails if the check-out date is later the same day", func(t *testing.T) {
		checkIn, checkOut := times("2019-Feb-03 at 9:00am", "2019-Feb-03 at 2:00pm")
		err := svc.NewBooker(mockHotels).Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("fails if check-out date is before check-in", func(t *testing.T) {
		checkIn, checkOut := times("2019-Feb-03", "2019-Feb-02")
		err := svc.NewBooker(mockHotels).Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("succeeds if the check-out date is one day after check-in", func(t *testing.T) {
		checkIn, checkOut := times("2019-Feb-03", "2019-Feb-04")
		err := svc.NewBooker(mockHotels).Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.NoError(t, err)
	})

	t.Run("succeeds if the check-out date is one month after checkin", func(t *testing.T) {
		checkIn, checkOut := times("2019-Feb-03", "2019-Mar-03")
		err := svc.NewBooker(mockHotels).Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.NoError(t, err)
	})

	t.Run("ignores the hours minutes component of the times", func(t *testing.T) {
		checkIn, checkOut := times("2019-Feb-03 at 11:59pm", "2019-Feb-04 at 12:01am")
		err := svc.NewBooker(mockHotels).Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
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

func TestBookAvailability(t *testing.T) {
	const alice = svc.EmployeeID("alice")
	const bob = svc.EmployeeID("bob")

	t.Run("fails if booking a room that is already booked", func(t *testing.T) {
		theHotel := svc.HotelID("the-hotel")
		checkIn, checkOut := times("2019-Mar-03", "2019-Mar-04")
		oneRoom := &HotelsMock{}
		oneRoom.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   theHotel,
			Name: "The Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "04"},
			},
		}
		booker := svc.NewBooker(oneRoom)

		err := booker.Book(alice, theHotel, svc.RoomSingle, checkIn, checkOut)
		assert.NoError(t, err)
		err = booker.Book(bob, theHotel, svc.RoomSingle, checkIn, checkOut)
		assert.Equal(t, svc.ErrNotAvailable, err)
	})

	t.Run("fails if two bookings partially overlap", func(t *testing.T) {
		theHotel := svc.HotelID("the-hotel")
		firstCheckIn, firstCheckOut := times("2019-Mar-03", "2019-Mar-22")
		secondCheckIn, secondCheckOut := times("2019-Mar-15", "2019-Apr-02")
		oneRoom := &HotelsMock{}
		oneRoom.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   theHotel,
			Name: "The Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "04"},
			},
		}
		booker := svc.NewBooker(oneRoom)

		err := booker.Book(alice, theHotel, svc.RoomSingle, firstCheckIn, firstCheckOut)
		assert.NoError(t, err)
		err = booker.Book(bob, theHotel, svc.RoomSingle, secondCheckIn, secondCheckOut)
		assert.Equal(t, svc.ErrNotAvailable, err)
	})

	t.Run("ignores bookings that do not overlap with the booking period", func(t *testing.T) {
		theHotel := svc.HotelID("the-hotel")
		marchCheckIn, marchCheckOut := times("2019-Mar-03", "2019-Mar-04")
		aprilCheckIn, aprilCheckOut := times("2019-Apr-22", "2019-Apr-28")
		oneRoom := &HotelsMock{}
		oneRoom.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   theHotel,
			Name: "The Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "04"},
			},
		}
		booker := svc.NewBooker(oneRoom)

		err := booker.Book(alice, theHotel, svc.RoomSingle, marchCheckIn, marchCheckOut)
		assert.NoError(t, err)
		err = booker.Book(bob, theHotel, svc.RoomSingle, aprilCheckIn, aprilCheckOut)
		assert.NoError(t, err)
	})

	t.Run("allows simulateous bookings if multiple rooms available", func(t *testing.T) {
		theHotel := svc.HotelID("the-hotel")
		checkIn, checkOut := times("2019-Mar-03", "2019-Mar-04")
		oneRoom := &HotelsMock{}
		oneRoom.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   theHotel,
			Name: "The Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "03"},
				svc.Room{Type: svc.RoomSingle, Number: "04"},
			},
		}
		booker := svc.NewBooker(oneRoom)

		err := booker.Book(alice, theHotel, svc.RoomSingle, checkIn, checkOut)
		assert.NoError(t, err)
		err = booker.Book(bob, theHotel, svc.RoomSingle, checkIn, checkOut)
		assert.NoError(t, err)
	})

	t.Run("fails simulaneous booking if only one available room of the requested type", func(t *testing.T) {
		theHotel := svc.HotelID("the-hotel")
		checkIn, checkOut := times("2019-Mar-03", "2019-Mar-04")
		oneRoom := &HotelsMock{}
		oneRoom.Func.GetHotel.Returns.Hotel = &svc.Hotel{
			ID:   theHotel,
			Name: "The Hotel",
			Rooms: []svc.Room{
				svc.Room{Type: svc.RoomSingle, Number: "03"},
				svc.Room{Type: svc.RoomDouble, Number: "04"},
			},
		}
		booker := svc.NewBooker(oneRoom)

		err := booker.Book(alice, theHotel, svc.RoomSingle, checkIn, checkOut)
		assert.NoError(t, err)
		err = booker.Book(bob, theHotel, svc.RoomSingle, checkIn, checkOut)
		assert.Equal(t, svc.ErrNotAvailable, err)
	})
}

func times(checkIn string, checkOut string) (checkInParsed time.Time, checkOutParsed time.Time) {
	return parseTime(checkIn), parseTime(checkOut)
}

func parseTime(str string) time.Time {
	var result time.Time
	var err error
	if strings.Contains(str, "at") {
		result, err = time.Parse(longForm, str)
	} else {
		result, err = time.Parse(shortForm, str)
	}
	if err != nil {
		panic(err)
	}
	return result
}
