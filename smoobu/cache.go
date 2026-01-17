package smoobu

import "strings"

type Cache struct {
	Apartments []Apartment
	User       *User
}

var Cached = &Cache{}

func (c *Cache) GetMainProperty() *Apartment {
	for _, apartment := range c.Apartments {
		if strings.Contains(strings.ToLower(apartment.Name), "clos du haslach") {
			return &apartment
		}
	}
	return nil
}

func (c *Cache) GetSubProperties() []Apartment {
	subs := []Apartment{}
	for _, apartment := range c.Apartments {
		if !strings.Contains(strings.ToLower(apartment.Name), "clos du haslach") {
			subs = append(subs, apartment)
		}
	}
	return subs
}

func (c *Cache) GetPropertyByName(name string) *Apartment {
	for _, apartment := range c.Apartments {
		if apartment.Name == name {
			return &apartment
		}
	}
	return nil
}
