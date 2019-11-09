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
const shortForm = "2006-Jan-02"
const longForm = "2006-Jan-02 at 3:04pm"

func TestBook(t *testing.T) {
	booker := svc.NewBooker()

	t.Run("returns an error if check-out date matches check-in", func(t *testing.T) {
		checkIn, _ := time.Parse(shortForm, "2019-Feb-03")
		checkOut := checkIn

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("returns an error if the check-out date is later the same day", func(t *testing.T) {
		checkIn, _ := time.Parse(longForm, "2019-Feb-03 at 9:00am")
		checkOut, _ := time.Parse(longForm, "2019-Feb-03 at 2:00pm")

		err := booker.Book(anyEmployee, anyHotel, anyRoom, checkIn, checkOut)
		assert.Equal(t, svc.ErrCheckoutInvalid, err)
	})

	t.Run("returns an error if check-out date is before check-in", func(t *testing.T) {
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
