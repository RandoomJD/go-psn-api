package psn

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

const trophyGroupSuffix = "-tpy."
const trophyGroupApi = "trophy/v1/trophyTitles"

// trophyURL returns https։//[region]-tpy․np․community․playstation․net
func trophyURL(region Region) *url.URL {
	new := baseURL.JoinPath()
	new.Host = string(region) + trophyGroupSuffix + new.Host
	return new
}

type TrophyTitle struct {
	Name            string        `json:"trophyTitleName"`
	Detail          string        `json:"trophyTitleDetail"`
	IconURL         string        `json:"trophyTitleIconUrl"`
	Platfrom        string        `json:"trophyTitlePlatfrom"` // typo in Sony's response
	DefinedTrophies TrophyCount   `json:"definedTrophies"`
	TrophyGroups    []TrophyGroup `json:"trophyGroups"`
}

type TrophyGroup struct {
	ID              string       `json:"trophyGroupId"`
	Name            string       `json:"trophyGroupName"`
	Detail          string       `json:"trophyGroupDetail"`
	IconURL         string       `json:"trophyGroupIconUrl"`
	SmallIconURL    string       `json:"trophyGroupSmallIconUrl"`
	DefinedTrophies TrophyCount  `json:"definedTrophies"`
	ComparedUser    UserProgress `json:"comparedUser"`
}

func (api *Api) baseTrophyQuery(username string, fields ...string) url.Values {
	var sb strings.Builder
	sb.WriteString(trophyFields)
	for _, f := range fields {
		sb.WriteString(",")
		sb.WriteString(f)
	}
	q := url.Values{}
	q.Add("fields", sb.String())
	q.Add("comparedUser", username)
	q.Add("npLanguage", string(api.lang))
	return q
}

// Method retrieves user's trophy groups
func (api *Api) GetTrophyGroups(ctx context.Context, trophyTitleId, username string) (TrophyTitle, error) {
	headers := api.trophyHeaders()

	url := trophyURL(api.region).JoinPath(trophyGroupApi, trophyTitleId, "trophyGroups")
	url.RawQuery = api.baseTrophyQuery(username).Encode()

	var resp TrophyTitle
	err := api.get(ctx, url, headers, &resp)
	if err != nil {
		return TrophyTitle{}, fmt.Errorf("can't do GET request: %w", err)
	}
	return resp, nil
}
