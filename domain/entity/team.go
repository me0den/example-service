package entity

// Team defines data model for resource Team struct.
type Team struct {
	ID    string `json:"id"`
	Owner string `json:"userID"`
}
