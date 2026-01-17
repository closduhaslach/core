package smoobu

import (
	"fmt"
	"net/url"
)

type GetApartmentResponse struct {
	Apartments []Apartment `json:"apartments"`
}

type CreateBookingErrorResponse struct {
	Status             int                 `json:"status"`
	Title              string              `json:"title"`
	Detail             string              `json:"detail"`
	ValidationMessages *ValidationMessages `json:"validation_messages,omitempty"`
}

type GetAvailabilityResponseEmpty struct {
	AvailableApartments []int `json:"availableApartments"`
	Prices              any   `json:"prices"`
	ErrorMessages       any   `json:"errorMessages"`
}

type GetAvailabilityResponse struct {
	AvailableApartments []int                               `json:"availableApartments"`
	Prices              map[string]Price                    `json:"prices"`
	ErrorMessages       map[string]AvailabilityErrorMessage `json:"errorMessages"`
}

type GetBookingsResponse struct {
	PageCount  int       `json:"page_count"`
	PageSize   int       `json:"page_size"`
	TotalItems int       `json:"total_items"`
	Page       int       `json:"page"`
	Bookings   []Booking `json:"bookings"`
}

func (b *GetBookingsResponse) NextPage(c Client) (*GetBookingsResponse, error) {
	if b.PageCount == b.Page {
		return nil, fmt.Errorf("no more pages")
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", b.Page+1))

	return c.GetBookings(params)
}

func (b *GetBookingsResponse) FilteredBookingsFunc(fn func(Booking) bool) []Booking {
	var filtered []Booking
	for _, booking := range b.Bookings {
		if fn(booking) {
			filtered = append(filtered, booking)
		}
	}
	return filtered
}

func (b *GetBookingsResponse) BlockedBookings() []Booking {
	return b.FilteredBookingsFunc(func(booking Booking) bool {
		return booking.IsBlockedBooking
	})
}

func (b *GetBookingsResponse) NonBlockedBookings() []Booking {
	return b.FilteredBookingsFunc(func(booking Booking) bool {
		return !booking.IsBlockedBooking
	})
}

func (b *GetBookingsResponse) BookingsByApartmentID(apartmentID int) []Booking {
	return b.FilteredBookingsFunc(func(booking Booking) bool {
		return booking.Apartment.ID == apartmentID
	})
}

func (b *GetBookingsResponse) BookingsByApartmentIDs(apartmentIDs []int) []Booking {
	idSet := make(map[int]struct{})
	for _, id := range apartmentIDs {
		idSet[id] = struct{}{}
	}
	return b.FilteredBookingsFunc(func(booking Booking) bool {
		_, exists := idSet[booking.Apartment.ID]
		return exists
	})
}
