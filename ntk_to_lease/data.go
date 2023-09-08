package ntk_to_lease

type Data struct {
	Operator  string  `json:"operator_name"`
	DateTime  string  `json:"datetime"`
	Amount    float32 `json:"amount"`
	Reference string  `json:"ref,omitempty"`
	Vrm       string  `json:"vrm,omitempty"`
	NTKUrl    string  `json:"ntk,omitempty"`
	AppealURL string  `json:"appeal_url,omitempty"`
	PayURL    string  `json:"pay_url,omitempty"`
}
