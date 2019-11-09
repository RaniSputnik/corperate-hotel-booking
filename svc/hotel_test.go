package svc_test

import (
	"testing"

	"github.com/RaniSputnik/corperate-hotel-booking/svc"
	"github.com/stretchr/testify/assert"
)

func TestHotelOffers(t *testing.T) {
	singleRooms := &svc.Hotel{
		ID:   "single-rooms",
		Name: "Single Rooms",
		Rooms: []svc.Room{
			svc.Room{Type: svc.RoomSingle, Number: "04"},
		},
	}

	doubleRooms := &svc.Hotel{
		ID:   "double-rooms",
		Name: "Double Rooms",
		Rooms: []svc.Room{
			svc.Room{Type: svc.RoomDouble, Number: "42"},
		},
	}

	suiteRooms := &svc.Hotel{
		ID:   "suite-rooms",
		Name: "Suite Rooms",
		Rooms: []svc.Room{
			svc.Room{Type: svc.RoomSuite, Number: "42"},
		},
	}

	allRooms := &svc.Hotel{
		ID:   "suite-rooms",
		Name: "Suite Rooms",
		Rooms: []svc.Room{
			svc.Room{Type: svc.RoomSingle, Number: "42"},
			svc.Room{Type: svc.RoomDouble, Number: "42"},
			svc.Room{Type: svc.RoomSuite, Number: "42"},
		},
	}

	t.Run("returns false if no rooms of the given type are offered", func(t *testing.T) {
		assert.False(t, singleRooms.Offers(svc.RoomDouble), "Double rooms should not be offered")
		assert.False(t, singleRooms.Offers(svc.RoomSuite), "Suite rooms should not be offered")

		assert.False(t, doubleRooms.Offers(svc.RoomSingle), "Single rooms should not be offered")
		assert.False(t, doubleRooms.Offers(svc.RoomSuite), "Suite rooms should not be offered")

		assert.False(t, suiteRooms.Offers(svc.RoomSingle), "Single rooms should not be offered")
		assert.False(t, suiteRooms.Offers(svc.RoomDouble), "Double rooms should not be offered")
	})

	t.Run("returns true if any rooms of the given type are offered", func(t *testing.T) {
		assert.True(t, singleRooms.Offers(svc.RoomSingle), "Single rooms should be offered")
		assert.True(t, doubleRooms.Offers(svc.RoomDouble), "Double rooms should be offered")
		assert.True(t, suiteRooms.Offers(svc.RoomSuite), "Suite rooms should be offered")

		assert.True(t, allRooms.Offers(svc.RoomSingle), "Single rooms should be offered")
		assert.True(t, allRooms.Offers(svc.RoomDouble), "Double rooms should be offered")
		assert.True(t, allRooms.Offers(svc.RoomSuite), "Suite rooms should be offered")
	})
}
