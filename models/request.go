package models

type BeginRegRequest struct {
	Username string `json:"username"`
}

type FinishRegRequest struct {
	Id    int    `json:"id"`
	RawID string `json:"rawId"`
	Type  string `json:"type"`

	Response struct {
		AttestationObject string `json:"attestationObject"`
		ClientDataJSON    string `json:"clientDataJSON"`
	}
}
