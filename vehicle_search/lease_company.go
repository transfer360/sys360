package vehicle_search

type LeaseCompany struct {
	LeaseCompanyId int    `json:"lease_company_id"`
	CompanyName    string `json:"company_name"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	AddressLine3   string `json:"address_line3"`
	City           string `json:"city"`
	PostCode       string `json:"postcode"`
}
