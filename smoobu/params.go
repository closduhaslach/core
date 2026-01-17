package smoobu

import (
	"net/url"
	"strconv"
)

func NewGetBookingsParams() *GetBookingsParams {
	return &GetBookingsParams{
		v: url.Values{},
	}
}

type GetBookingsParams struct {
	v url.Values
}

func (g *GetBookingsParams) From(from string) *GetBookingsParams {
	g.v.Add("from", from)
	return g
}

func (g *GetBookingsParams) To(to string) *GetBookingsParams {
	g.v.Add("to", to)
	return g
}

func (g *GetBookingsParams) ApartmentID(id int) *GetBookingsParams {
	g.v.Add("apartmentId", strconv.Itoa(id))
	return g
}

func (g *GetBookingsParams) ExcludeBlocked(exclude bool) *GetBookingsParams {
	g.v.Add("excludeBlocked", strconv.FormatBool(exclude))
	return g
}

func (g *GetBookingsParams) Values() url.Values {
	return g.v
}
