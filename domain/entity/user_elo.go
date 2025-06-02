package entity

const (
	DefaultElo = 1000
)

type UserElo struct {
	UserID string `json:"userID"`
	Elo    int    `json:"elo"`
}

func (e *UserElo) Clone() *UserElo {
	return &UserElo{
		UserID: e.UserID,
		Elo:    e.Elo,
	}
}

func NewUserDefaultElo(userID string) *UserElo {
	return &UserElo{
		UserID: userID,
		Elo:    DefaultElo,
	}
}
