package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APICardImage struct {
	ImageURL string `json:"image_url"`
}

type APICard struct {
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Type              string         `json:"type"`
	HumanReadableType string         `json:"humanReadableCardType"`
	FrameType         string         `json:"frameType"`
	Desc              string         `json:"desc"`
	Race              string         `json:"race"`
	Atk               int            `json:"atk"`
	Def               int            `json:"def"`
	Level             int            `json:"level"`
	Attribute         string         `json:"attribute"`
	Archetype         string         `json:"archetype"`
	CardImages        []APICardImage `json:"card_images"`
	LinkValue         int            `json:"linkval"`
	LinkMarkers       []string       `json:"linkmarkers"`
	Scale             int            `json:"scale"`
	ImageURL          string
}

type CardResponse struct {
	Data []APICard `json:"data"`
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

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw Body:\n", string(bodyBytes))

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var result CardResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Data) == 0 || len(result.Data[0].CardImages) == 0 {
		return nil, errors.New("card not found or no card image found")
	}

	card := result.Data[0]
	card.ImageURL = card.CardImages[0].ImageURL
	return &card, nil
}

func FetchRandomCards(n int) ([]APICard, error) {
	url := "https://db.ygoprodeck.com/api/v7/randomcard.php"
	var cards []APICard

	for i := 0; i < n; i++ {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch card %d: %w", i, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}

		var apiResp CardResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("failed to decode card %d: %w", i, err)
		}

		if len(apiResp.Data) == 0 || apiResp.Data[0].ID == 0 {
			fmt.Println("Invalid card received")
			continue
		}

		cards = append(cards, apiResp.Data[0])
	}

	return cards, nil
}
