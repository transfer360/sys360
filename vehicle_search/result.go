package vehicle_search

type SearchResult struct {
	Sref           string        `json:"sref"`
	Request        SearchRequest `json:"request"`
	IsLeaseVehicle bool          `json:"is_lease_vehicle"`
	Response       LeaseCompany  `json:"response"`
	Processed      string        `json:"processed"`
}
