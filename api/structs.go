package api

// Register models the register operation input
type Register struct {
	Key            *string `json:"key"`
	Url            *string `json:"url"`
	HealthEndpoint *string `json:"healthEndpoint"`
	TTL            *int    `json:"ttl"`
}
