package psn

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const profileSuffix = "-prof."

func profileURL(region Region) *url.URL {
	new := baseURL.JoinPath()
	new.Host = string(region) + profileSuffix + new.Host
	return new
}

type AvatarUrls struct {
	Size      string `json:"size"`
	AvatarURL string `json:"avatarUrl"`
}
type TrophyCount struct {
	Platinum int `json:"platinum"`
	Gold     int `json:"gold"`
	Silver   int `json:"silver"`
	Bronze   int `json:"bronze"`
}
type TrophySummary struct {
	Level          int         `json:"level"`
	Progress       int         `json:"progress"`
	EarnedTrophies TrophyCount `json:"earnedTrophies"`
}

type PersonalDetail struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Presences struct {
}

type ConsoleAvailability struct {
	AvailabilityStatus string `json:"availabilityStatus"`
}

type User struct {
	OnlineID string `json:"onlineId"`
}
type Profile struct {
	User
	NpID                  string              `json:"npId"`
	AvatarUrls            []AvatarUrls        `json:"avatarUrls"`
	Plus                  int                 `json:"plus"`
	AboutMe               string              `json:"aboutMe"`
	LanguagesUsed         []string            `json:"languagesUsed"`
	TrophySummary         TrophySummary       `json:"trophySummary"`
	IsOfficiallyVerified  bool                `json:"isOfficiallyVerified"`
	PersonalDetail        PersonalDetail      `json:"personalDetail"`
	PersonalDetailSharing string              `json:"personalDetailSharing"`
	PrimaryOnlineStatus   string              `json:"primaryOnlineStatus"`
	Presences             []Presences         `json:"presences"`
	FriendRelation        string              `json:"friendRelation"`
	Blocking              bool                `json:"blocking"`
	MutualFriendsCount    int                 `json:"mutualFriendsCount"`
	Following             bool                `json:"following"`
	FollowerCount         int                 `json:"followerCount"`
	ConsoleAvailability   ConsoleAvailability `json:"consoleAvailability"`
}

var profileFields = "onlineId,aboutMe,consoleAvailability,languagesUsed,avatarUrls,personalDetail,personalDetail(@default,profilePictureUrls),primaryOnlineStatus,trophySummary(level,progress,earnedTrophies),plus,isOfficiallyVerified,friendRelation,personalDetailSharing,presences(@default,platform),npId,blocking,following,currentOnlineId,displayableOldOnlineId,mutualFriendsCount,followerCount"

func baseProfileQuery() url.Values {
	query := url.Values{}
	query.Add("fields", profileFields)
	query.Add("profilePictureSizes", "s,m,l")
	query.Add("avatarSizes", "s,m,l")
	query.Add("languagesUsedLanguageSet", "set4")
	return query
}

func (api *AuthedApi) profileHeaders() http.Header {
	var h http.Header
	h.Add("authorization", fmt.Sprintf("Bearer %s", api.tokens.Access))
	return h
}

// Method retrieves user profile info by PSN id
func (api *AuthedApi) GetProfileRequest(ctx context.Context, username string) (profile *Profile, err error) {
	type Response struct {
		Profile Profile `json:"profile"`
	}

	headers := api.profileHeaders()

	url := profileURL(api.region).JoinPath(username, "profile2")

	url.RawQuery = baseProfileQuery().Encode()

	var resp Response
	err = api.get(ctx, url, headers, &resp)
	if err != nil {
		return nil, fmt.Errorf("can't do GET request: %w", err)
	}
	return &resp.Profile, nil
}
