package twillio

type Twillio struct {
	AccountSID         string
	AuthToken          string
	TwillioPhoneNumber string
	NgrokUrl           string
}

type InitiateCallResponse struct {
	RecipientNumber string `json:"RecipientNumber"`
	TargetName      string `json:"TargetName"`
	Service         string `json:"Service"`
	LenDigits       int    `json:"LenDigits"`
}

type SidResponse struct {
	Sid string `json:"Sid"`
}
