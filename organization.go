package beasiswakita

import "time"

type Organization struct {
	ID                int       `db:"id" json:"id"`
	Name              string    `db:"name" json:"name"`
	Position          string    `db:"position" json:"position"`
	OrganizationName  string    `db:"organization_name" json:"organization_name"`
	OrganizationEmail string    `db:"organization_email" json:"organization_email"`
	Address           string    `db:"address" json:"address"`
	City              string    `db:"city" json:"city"`
	Region            string    `db:"region" json:"region"`
	Country           string    `db:"country" json:"country"`
	Zipcode           *string   `db:"zipcode" json:"zipcode,omitempty"`
	Website           *string   `db:"website" json:"website,omitempty"`
	Logo              *string   `db:"logo" json:"logo,omitempty"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	UserID            int       `db:"user_id" json:"-"`
}

func (o *Organization) Parse(data map[string]interface{}) error {
	o.Name = data["name"].(string)
	o.Position = data["position"].(string)
	o.OrganizationName = data["organization_name"].(string)
	o.OrganizationEmail = data["organization_email"].(string)
	o.Address = data["address"].(string)
	o.City = data["city"].(string)
	o.Region = data["region"].(string)
	o.Country = data["country"].(string)

	zipcode := data["zipcode"].(string)
	website := data["website"].(string)
	logo := data["logo"].(string)

	o.Zipcode = &zipcode
	o.Website = &website
	o.Logo = &logo

	return nil
}

func (o *Organization) Map() map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = o.ID
	data["name"] = o.Name
	data["position"] = o.Position
	data["organization_name"] = o.OrganizationName
	data["organization_email"] = o.OrganizationEmail
	data["address"] = o.Address
	data["city"] = o.City
	data["region"] = o.Region
	data["country"] = o.Country
	data["zipcode"] = o.Zipcode
	data["website"] = o.Website
	data["logo"] = o.Logo

	return data
}
