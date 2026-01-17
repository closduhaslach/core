// Package smoobu
package smoobu

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
}

type Price struct {
	Price         int            `json:"price"`
	Currency      string         `json:"currency"`
	PriceElements []PriceElement `json:"priceElements,omitempty"`
}

type AvailabilityErrorMessage struct {
	ErrorCode                    int      `json:"errorCode"`
	Message                      string   `json:"message"`
	MinimumLengthOfStay          int      `json:"minimumLengthOfStay,omitempty"`
	NumberOfGuests               int      `json:"numberOfGuests,omitempty"`
	LeadTime                     int      `json:"leadTime,omitempty"`
	MinimumLengthBetweenBookings int      `json:"minimumLengthBetweenBookings,omitempty"`
	ArrivalDays                  []string `json:"arrivalDays,omitempty"`
}

type Event[T any] struct {
	Action string `json:"action"`
	User   int    `json:"user"`
	Data   T      `json:"data"`
}

type Booking struct {
	ID                 int            `json:"id"`
	ReferenceID        string         `json:"reference-id,omitempty"`
	Type               string         `json:"type"`
	Arrival            string         `json:"arrival"`
	Departure          string         `json:"departure"`
	CreatedAt          string         `json:"created-at"`
	ModifiedAt         string         `json:"modifiedAt"`
	Apartment          Apartment      `json:"apartment"`
	Channel            Channel        `json:"channel"`
	GuestName          string         `json:"guest-name"`
	Firstname          string         `json:"firstname"`
	Lastname           string         `json:"lastname"`
	Email              string         `json:"email"`
	Phone              *string        `json:"phone"` // nullable
	Adults             int            `json:"adults"`
	Children           int            `json:"children"`
	CheckIn            string         `json:"check-in"`
	CheckOut           string         `json:"check-out"`
	Notice             string         `json:"notice"`
	AssistantNotice    string         `json:"assistant-notice"`
	Price              float64        `json:"price"`
	PriceDetails       string         `json:"price-details"`
	CityTax            *float64       `json:"city-tax"` // nullable
	PricePaid          string         `json:"price-paid"`
	CommissionIncluded *float64       `json:"commission-included"` // nullable
	Prepayment         float64        `json:"prepayment"`
	PrepaymentPaid     string         `json:"prepayment-paid"`
	Deposit            *float64       `json:"deposit"` // nullable
	DepositPaid        string         `json:"deposit-paid"`
	Language           string         `json:"language"`
	GuestAppURL        string         `json:"guest-app-url"`
	IsBlockedBooking   bool           `json:"is-blocked-booking"`
	GuestID            int            `json:"guestId"`
	Related            []Apartment    `json:"related"`
	PriceElements      []PriceElement `json:"priceElements,omitempty"`
}

func (b *Booking) IsBlockedBookingManaged() bool {
	return b.IsBlockedBooking && b.Firstname == "Blocker" && b.Lastname == "Smoobuwebhook"
}

func (b *Booking) IsBlockedBookingUnmanaged() bool {
	return b.IsBlockedBooking && (b.Firstname != "Blocker" || b.Lastname != "Smoobuwebhook")
}

type Apartment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Channel struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channel_id"`
	Name      string `json:"name"`
}

type PriceElement struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Amount       float64  `json:"amount"`
	Quantity     *int     `json:"quantity"` // nullable
	Tax          *float64 `json:"tax"`      // nullable
	CurrencyCode string   `json:"currencyCode"`
	SortOrder    int      `json:"sortOrder"`
}

type Address struct {
	Street     string `json:"street,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Location   string `json:"location,omitempty"`
}

type ValidationMessages struct {
	Error string `json:"error"`
}

type Blocker struct {
	Arrival   string
	Departure string
	Apartment Apartment
}
