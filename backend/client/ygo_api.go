package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type APICard struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	HumanReadableType string `json:"humanReadableCardType"`
	FrameType         string `json:"frameType"`
	Desc              string `json:"desc"`
	Race              string `json:"race"`
	Atk               int    `json:"atk"`
	Def               int    `json:"def"`
	Level             int    `json:"level"`
	Attribute         string `json:"attribute"`
	Archetype         string `json:"archetype"`
	ImageURL          string
	LinkValue         int      `json:"linkval"`
	LinkMarkers       []string `json:"linkmarkers"`
	Scale             int      `json:"scale"`
}

func FetchCardByIDOrName(id int, name string) (*APICard, error) {
	var endpoint string
	if id > 0 {
		endpoint = fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?id=%d", id)
	} else {
		endpoint = fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?name=%s", url.QueryEscape(name))
	}

	resp, err := http.Get(endpoint)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("external API error")
	}
	defer resp.Body.Close()

	var result struct {
		Data []APICard `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, errors.New("card not found")
	}

	return &result.Data[0], nil
}
