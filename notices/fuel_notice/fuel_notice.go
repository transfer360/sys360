package fuel_notice

type Data struct {
	Source                string `json:"source" firestore:"-"`
	DocumentID            string `json:"-" firestore:"-"`
	SearchReference       string `json:"sref" firestore:"sref"`
	OperatorName          string `json:"operator_name" firestore:"operator_name"`
	NoticeType            string `json:"-" firestore:"notice_type"`
	NoticeTkn             string `json:"notice_tkn" firestore:"notice_tkn"`
	NoticeNumber          string `json:"notice_number" firestore:"notice_number"`
	VehicleRegistration   string `json:"vehicle_registration" firestore:"vehicle_registration"`
	Contravention         string `json:"contravention" firestore:"contravention"`
	ContraventionDateTime string `json:"contravention_datetime" firestore:"contravention_datetime"`
	EntryExit             struct {
		Entry string `json:"entry,omitempty" firestore:"entry,omitempty"`
		Exit  string `json:"exit,omitempty" firestore:"exit,omitempty"`
	} `json:"entry_exit_datetime,omitempty" firestore:"entry_exit_datetime,omitempty"`
	//Observation struct {
	//	From string `json:"from,omitempty" firestore:"from,omitempty"`
	//	To   string `json:"to,omitempty" firestore:"to,omitempty"`
	//} `json:"observation_datetime,omitempty" firestore:"observation_datetime"`
	Location       string `json:"location" firestore:"location"`
	NoticeToKeeper struct {
		File     string `json:"file" firestore:"-"`
		URL      string `json:"url" firestore:"url"`
		Received string `json:"-" firestore:"received"`
	} `json:"notice_to_keeper" firestore:"notice_to_keeper"`
	//Letter struct {
	//	File     string `json:"file" firestore:"-"`
	//	URL      string `json:"url" firestore:"url"`
	//	Received string `json:"-" firestore:"received"`
	//	Type     int    `json:"type" firestore:"type"`
	//} `json:"letter" firestore:"letter"`
	//Pofa              bool     `json:"pofa" firestore:"pofa"`
	TotalDue          float64  `json:"total_due" firestore:"total_due"`
	ReducedAmount     float64  `json:"reduced_amount" firestore:"reduced_amount"`
	ReducedPeriodEnds string   `json:"reduce_period_ends" firestore:"reduce_period_ends"`
	Photos            []string `json:"photos,omitempty" firestore:"photos,omitempty"`
	PaymentURL        string   `json:"payment_url,omitempty" firestore:"payment_url,omitempty"`
	AppealURL         string   `json:"appeal_url,omitempty" firestore:"appeal_url,omitempty"`
	IssuerID          string   `json:"-" firestore:"-"`
	FleetID           int      `json:"-" firestore:"-"`
	System            struct {
		IgnoreNoticeExistsCheck bool /* Ignore the has already been uploaded check */
		NoticeUpload            struct {
			ResendNotice      bool
			NoticeDocumentID  string
			NoticeTkn         string
			SaveSearchDocment bool
		}
	}
}
