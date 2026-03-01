package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"steam-bubbles/internal/domain"
)

const ownedGamesEndpoint = "https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/"

type apiGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

type ownedGamesResponse struct {
	Response struct {
		GameCount int       `json:"game_count"`
		Games     []apiGame `json:"games"`
	} `json:"response"`
}

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) GetOwnedGames(_ context.Context, steamID string) ([]domain.Game, error) {
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("steamid", steamID)
	params.Set("include_appinfo", "true")
	params.Set("include_played_free_games", "true")
	params.Set("format", "json")

	resp, err := c.http.Get(ownedGamesEndpoint + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("steam: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam: API returned status %d", resp.StatusCode)
	}

	var raw ownedGamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("steam: decode error: %w", err)
	}

	if len(raw.Response.Games) == 0 {
		return nil, fmt.Errorf("steam: no games returned — ensure your Steam profile is set to public")
	}

	games := make([]domain.Game, len(raw.Response.Games))
	for i, g := range raw.Response.Games {
		games[i] = domain.Game{
			AppID:           g.AppID,
			Name:            g.Name,
			PlaytimeMinutes: g.PlaytimeForever,
		}
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].PlaytimeMinutes > games[j].PlaytimeMinutes
	})

	return games, nil
}
