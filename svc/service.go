package svc

type HotelID string
type EmployeeID string

const (
	AllEmployees = EmployeeID("")
)

type CompanyID string

type RoomType string

const (
	RoomSingle = RoomType("room:single")
	RoomDouble = RoomType("room:double")
	RoomSuite  = RoomType("room:suite")
)

type Room struct {
	Type   RoomType
	Number string
}

type Hotel struct {
	ID    HotelID
	Name  string
	Rooms []Room
}

type Hotels interface {
	AddHotel(id HotelID, name string, rooms []Room) (*Hotel, error)
}

type Employees interface {
	AddEmployee(id EmployeeID, company CompanyID)
	DeleteEmployee(id EmployeeID)
}

type Policies interface {
	// AddPolicy adds a room booking policy for a given company
	// and optionally, a specific employee. Use the "AllEmployees"
	// constant to apply the policy to all employees within a company.
	//
	// Specifying room types narrows a policy to those particular
	// room types specified. Providing no room types applys the policy
	// to all room types.
	AddPolicy(company CompanyID, employee EmployeeID, rooms ...RoomType)

	// Allow returns whether or not a given employee may book the
	// given room type.
	Allow(employee EmployeeID, room RoomType) bool
}
