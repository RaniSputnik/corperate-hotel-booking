package svc_test

import (
	"errors"

	"github.com/RaniSputnik/corperate-hotel-booking/svc"
)

type HotelsMock struct {
	Func struct {
		GetHotel struct {
			CalledWith struct {
				ID svc.HotelID
			}
			Returns struct {
				Hotel *svc.Hotel
				Err   error
			}
		}
	}
}

func (m *HotelsMock) AddHotel(id svc.HotelID, name string, rooms []svc.Room) (*svc.Hotel, error) {
	return nil, errors.New("not implemented")
}

func (m *HotelsMock) GetHotel(id svc.HotelID) (*svc.Hotel, error) {
	m.Func.GetHotel.CalledWith.ID = id
	return m.Func.GetHotel.Returns.Hotel, m.Func.GetHotel.Returns.Err
}
