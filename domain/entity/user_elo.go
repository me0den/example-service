package entity

const (
	DefaultElo = 1000
)

// UserElo defines data model for resource UserElo struct.
type UserElo struct {
	UserID string `json:"userID"`
	Elo    int    `json:"elo"`
}

// Clone create a new object UserElo with exists value.
func (e *UserElo) Clone() *UserElo {
	return &UserElo{
		UserID: e.UserID,
		Elo:    e.Elo,
	}
}

// NewUserDefaultElo create a new object UserElo with default value.
func NewUserDefaultElo(userID string) *UserElo {
	return &UserElo{
		UserID: userID,
		Elo:    DefaultElo,
	}
}
