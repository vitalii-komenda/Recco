package domain

import (
	"context"
	"fmt"
)

type Game struct {
	AppID           int
	Name            string
	PlaytimeMinutes int
}

func (g Game) FormattedPlaytime() string {
	if g.PlaytimeMinutes == 0 {
		return "never played"
	}
	if g.PlaytimeMinutes < 60 {
		return fmt.Sprintf("%dm", g.PlaytimeMinutes)
	}
	h := g.PlaytimeMinutes / 60
	m := g.PlaytimeMinutes % 60
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

type Recommendation struct {
	Name     string `json:"name"`
	Reason   string `json:"reason"`
	Platform string `json:"platform"`
}

type GameRepository interface {
	GetOwnedGames(ctx context.Context, steamID string) ([]Game, error)
}

type Recommender interface {
	Recommend(ctx context.Context, games []Game) ([]Recommendation, error)
}
