package psn

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const trophiesApi = "trophy/v1/trophyTitles/"

type UserProgress struct {
	User
	Progress       int         `json:"progress"`
	EarnedTrophies TrophyCount `json:"earnedTrophies"`
	HiddenFlag     bool        `json:"hiddenFlag"`
	LastUpdateDate time.Time   `json:"lastUpdateDate"`
}

type UserEarned struct {
	User
	Earned bool `json:"earned"`
}

type Trophy struct {
	ID           int        `json:"trophyId"`
	Hidden       bool       `json:"trophyHidden"`
	Type         string     `json:"trophyType"`
	Name         string     `json:"trophyName"`
	Detail       string     `json:"trophyDetail"`
	IconURL      string     `json:"trophyIconUrl"`
	SmallIconURL string     `json:"trophySmallIconUrl"`
	Rare         int        `json:"trophyRare"`
	EarnedRate   string     `json:"trophyEarnedRate"`
	FromUser     UserEarned `json:"fromUser"`
}

func (api *Api) trophyHeaders() http.Header {
	h := api.profileHeaders()
	h.Add("Accept", "*/*")
	h.Add("Accept-Encoding", "gzip, deflate, br")
	return h
}

var trophyFields = "@default, trophySmallIconUrl"

// Method retrieves user's trophies
func (api *Api) GetTrophies(ctx context.Context, titleId, trophyGroupId, username string) ([]Trophy, error) {
	type Response struct {
		Trophies []Trophy `json:"trophies"`
	}

	headers := api.trophyHeaders()

	url := trophyURL(api.region).JoinPath(trophiesApi, titleId, "trophyGroups", trophyGroupId, "trophies")
	q := api.baseTrophyQuery(username, "trophyRare", "trophyEarnedRate")
	q.Add("visibleType", "1")
	url.RawQuery = q.Encode()

	var resp Response
	err := api.get(ctx, url, headers, &resp)
	if err != nil {
		return nil, fmt.Errorf("can't do GET request: %w", err)
	}
	return resp.Trophies, nil
}
