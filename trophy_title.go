package psn

import (
	"context"
	"fmt"
	"strconv"
)

type TrophyTitleFull struct {
	TrophyTitle
	NpCommunicationID       string       `json:"npCommunicationId"`
	TrophyTitleSmallIconURL string       `json:"trophyTitleSmallIconUrl"`
	HasTrophyGroups         bool         `json:"hasTrophyGroups"`
	ComparedUser            UserProgress `json:"comparedUser"`
	FromUser                UserProgress `json:"fromUser"`
}

var platforms = "PS3,PS4,PSVITA"

// GetTrophyTitles retrieves a user's trophy titles
func (api *AuthedApi) GetTrophyTitles(ctx context.Context, username string, limit, offset int) ([]TrophyTitleFull, error) {
	type Response struct {
		TotalResults int               `json:"totalResults"`
		Offset       int               `json:"offset"`
		Limit        int               `json:"limit"`
		TrophyTitles []TrophyTitleFull `json:"trophyTitles"`
	}

	headers := api.trophyHeaders()

	url := trophyURL(api.region)
	q := api.baseTrophyQuery(username)
	q.Add("platform", platforms)
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	url.RawQuery = q.Encode()

	var resp Response
	err := api.get(ctx, url, headers, &resp)
	if err != nil {
		return nil, fmt.Errorf("request for %s's trophy titles failed: %w", username, err)
	}
	return resp.TrophyTitles, nil
}
