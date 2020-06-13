package api

// Register models the register operation input
type Register struct {
	Key       *string `json:"key"`
	URL       *string `json:"url"`
	HealthURL *string `json:"healthUrl"`
}
