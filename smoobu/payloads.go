package smoobu

type BookingPayload struct {
	ArrivalDate      string         `json:"arrivalDate"`
	DepartureDate    string         `json:"departureDate"`
	ChannelID        int            `json:"channelId"`
	ApartmentID      int            `json:"apartmentId"`
	ArrivalTime      string         `json:"arrivalTime,omitempty"`
	DepartureTime    string         `json:"departureTime,omitempty"`
	FirstName        string         `json:"firstName,omitempty"`
	LastName         string         `json:"lastName,omitempty"`
	Notice           string         `json:"notice,omitempty"`
	Adults           int            `json:"adults,omitempty"`
	Children         int            `json:"children,omitempty"`
	Price            float64        `json:"price,omitempty"`
	PriceStatus      int            `json:"priceStatus,omitempty"`
	Prepayment       float64        `json:"prepayment,omitempty"`
	PrepaymentStatus int            `json:"prepaymentStatus,omitempty"`
	Deposit          float64        `json:"deposit,omitempty"`
	DepositStatus    int            `json:"depositStatus,omitempty"`
	Address          Address        `json:"address"`
	Country          string         `json:"country,omitempty"`
	Email            string         `json:"email,omitempty"`
	Phone            string         `json:"phone,omitempty"`
	Language         string         `json:"language,omitempty"`
	PriceElements    []PriceElement `json:"priceElements,omitempty"`
}

type AvailabilityPayload struct {
	ArrivalDate   string `json:"arrivalDate"`
	DepartureDate string `json:"departureDate"`
	Apartments    []int  `json:"apartments"`
	CustomerID    int    `json:"customerId"`
	Guests        int    `json:"guests,omitempty"`
	DiscountCode  string `json:"discountCode,omitempty"`
}

type UpdateBookingPayload struct {
	DepartureTime    string  `json:"departureTime"`
	ArrivalTime      string  `json:"arrivalTime"`
	Price            float64 `json:"price"`
	PriceStatus      int     `json:"priceStatus"`
	Prepayment       float64 `json:"prepayment"`
	PrepaymentStatus int     `json:"prepaymentStatus"`
	Deposit          float64 `json:"deposit"`
	DepositStatus    int     `json:"depositStatus"`
	Notice           string  `json:"notice"`
	AssistantNotice  string  `json:"assistantNotice"`
	GuestName        string  `json:"guestName"`
	GuestEmail       string  `json:"guestEmail"`
	GuestPhone       string  `json:"guestPhone"`
	Adults           int     `json:"adults"`
	Children         int     `json:"children"`
	Language         string  `json:"language"`
}
