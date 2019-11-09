package svc

type Hotel struct {
	ID    HotelID
	Name  string
	Rooms []Room
}

// Offers returns true or false whether the hotel contains
// any rooms of the given type.
func (h *Hotel) Offers(room RoomType) bool {
	for _, test := range h.Rooms {
		if test.Type == room {
			return true
		}
	}
	return false
}

// Options returns the rooms of the given room type
func (h *Hotel) Options(room RoomType) []Room {
	results := []Room{}
	for _, test := range h.Rooms {
		if test.Type == room {
			results = append(results, test)
		}
	}
	return results
}