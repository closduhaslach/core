package smoobu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetApartments() ([]Apartment, error) {
	resp, err := c.Get("/api/apartments")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get apartments: status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var apartments GetApartmentResponse
	err = json.NewDecoder(resp.Body).Decode(&apartments)
	if err != nil {
		return nil, err
	}
	return apartments.Apartments, nil
}

func (c *Client) GetBookings(params ...url.Values) (*GetBookingsResponse, error) {
	defaults := url.Values{}
	defaults.Set("page", "1")
	defaults.Set("pageSize", "10")
	defaults.Set("excludeBlocked", "false")
	defaults.Set("showCancellation", "false")

	params = append([]url.Values{defaults}, params...)

	allParams := url.Values{}
	for _, p := range params {
		for key, values := range p {
			for _, value := range values {
				allParams.Add(key, value)
			}
		}
	}

	slog.Debug("Fetching bookings")
	resp, err := c.Get(fmt.Sprintf("/api/reservations?%s", allParams.Encode()))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get bookings: status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	var bookings GetBookingsResponse
	err = json.NewDecoder(resp.Body).Decode(&bookings)
	if err != nil {
		return nil, err
	}
	return &bookings, nil
}

func (c *Client) GetBooking(bookingID int) (*Booking, error) {
	slog.Debug("Fetching booking", "booking_id", bookingID)
	resp, err := c.Get(fmt.Sprintf("/api/reservations/%d", bookingID))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get booking: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var booking Booking
	err = json.NewDecoder(resp.Body).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (c *Client) CreateBooking(booking BookingPayload) error {
	b, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	slog.Debug("Creating booking", "booking", string(b))
	resp, err := c.Post("/api/reservations", "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to create booking: status code %d", resp.StatusCode)
		errMsg, _ := io.ReadAll(resp.Body)
		slog.Error("Create booking failed", "error", err, "data", string(errMsg))
		return fmt.Errorf("failed to create booking: status code %d", resp.StatusCode)
	}

	var bookingResp CreateBookingErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&bookingResp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteBooking(bookingID int) error {
	slog.Debug("Deleting booking", "booking_id", bookingID)

	resp, err := c.Delete(fmt.Sprintf("/api/reservations/%d", bookingID))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete booking: status code %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteBookingSafe(bookingID int) error {
	slog.Debug("Deleting booking", "booking_id", bookingID)
	booking, err := c.GetBooking(bookingID)
	if err != nil {
		return err
	}

	if !booking.IsBlockedBooking {
		return fmt.Errorf("refusing to delete booking %d: not a blocker", bookingID)
	}

	return c.DeleteBooking(bookingID)
}

func (c *Client) GetUser() (*User, error) {
	resp, err := c.Get("/api/me")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) parseAvailResponse(resp *http.Response) (*GetAvailabilityResponse, error) {
	defer resp.Body.Close()

	var availability GetAvailabilityResponse
	var availabilityEmpty GetAvailabilityResponseEmpty

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	err := json.NewDecoder(tee).Decode(&availability)
	if err != nil {
		err = json.NewDecoder(&buf).Decode(&availabilityEmpty)
		if err != nil {
			return nil, err
		}
		return &GetAvailabilityResponse{
			AvailableApartments: availabilityEmpty.AvailableApartments,
			Prices:              map[string]Price{},
			ErrorMessages:       map[string]AvailabilityErrorMessage{},
		}, nil
	}
	return &availability, nil
}

func (c *Client) GetAvailability(payload AvailabilityPayload) (*GetAvailabilityResponse, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	slog.Debug("Checking availability", "payload", string(b))
	resp, err := c.Post("/booking/checkApartmentAvailability", "application/json", strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var availability *GetAvailabilityResponse
	availability, err = c.parseAvailResponse(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return availability, fmt.Errorf("failed to get availability: status code %d", resp.StatusCode)
	}
	return availability, nil
}

func (c *Client) UpdateBooking(bookingID int, payload UpdateBookingPayload) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	slog.Debug("Updating booking", "booking_id", bookingID, "payload", string(b))
	resp, err := c.Post(fmt.Sprintf("/api/reservations/%d", bookingID), "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errMsg, _ := io.ReadAll(resp.Body)
		slog.Error("Update booking failed", "data", string(errMsg))
		return fmt.Errorf("failed to update booking: status code %d", resp.StatusCode)
	}

	return nil
}

// BlockDates creates blocker bookings for the given date range for a specific apartment.
func (c *Client) BlockDates(blocker Blocker) error {
	slog.Info(
		"Blocking dates",
		"apartment_name", blocker.Apartment.Name,
		"apartment_id", blocker.Apartment.ID,
		"arrival", blocker.Arrival,
		"departure", blocker.Departure,
		"firstName", "Blocker",
		"lastName", "SmoobuWebhook",
		"email", "contact@closduhaslach.com",
		"phone", "+33388957345",
		"country", "FR",
		"address", "blocked",
	)
	return c.CreateBooking(BookingPayload{
		ArrivalDate:      blocker.Arrival,
		DepartureDate:    blocker.Departure,
		ChannelID:        11,
		ApartmentID:      blocker.Apartment.ID,
		FirstName:        "Blocker",
		LastName:         "SmoobuWebhook",
		Email:            "contact@closduhaslach.com",
		Price:            0,
		Prepayment:       0,
		Deposit:          0,
		PriceStatus:      1,
		PrepaymentStatus: 1,
		DepositStatus:    1,
		Country:          "FR",
		Address:          Address{Street: "blocked", PostalCode: "67000", Location: "blocked"},
	})
}

// UnblockDates removes blocker bookings for the given date range.
func (c *Client) UnblockDates(blocker Blocker) error {
	blockers, err := c.GetBookings(
		NewGetBookingsParams().
			From(blocker.Arrival).
			To(blocker.Departure).
			ApartmentID(blocker.Apartment.ID).
			Values(),
	)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, b := range blockers.BlockedBookings() {
		if !b.IsBlockedBooking {
			slog.Info("Skipping non-blocker booking", "bookingID", b.ID)
			continue
		}

		if b.Channel.ChannelID != 11 || b.Firstname != "Blocker" || b.Lastname != "Smoobuwebhook" {
			errs = append(errs, fmt.Errorf("could not remove blocker %d as it is not a managed blocker", b.ID))
			continue
		}

		if b.Arrival != blocker.Arrival || b.Departure != blocker.Departure {
			errs = append(errs, fmt.Errorf("could not remove blocker %d as it does not match exact date range: %s - %s", b.ID, b.Arrival, b.Departure))
			continue
		}

		slog.Info("Removing blocker", "bookingID", b.ID, "apartmentID", b.Apartment.ID, "apartmentName", b.Apartment.Name)
		err = c.DeleteBookingSafe(b.ID)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
