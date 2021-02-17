package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

type twitterUser struct {
	ID              uint64      `json:"id"`
	IDStr           string      `json:"id_str"`
	Name            string      `json:"name"`
	ScreenName      string      `db:"screen_name" json:"screen_name"`
	Location        string      `json:"location"`
	ProfileLocation interface{} `json:"profile_location"`
	Description     string      `json:"description"`
	URL             interface{} `json:"url"`
	Entities        struct {
		Description struct {
			Urls []interface{} `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	Protected       bool        `json:"protected"`
	FollowersCount  int         `json:"followers_count"`
	FriendsCount    int         `json:"friends_count"`
	ListedCount     int         `json:"listed_count"`
	CreatedAt       string      `json:"created_at"`
	FavouritesCount int         `json:"favourites_count"`
	UtcOffset       interface{} `json:"utc_offset"`
	TimeZone        interface{} `json:"time_zone"`
	GeoEnabled      bool        `json:"geo_enabled"`
	Verified        bool        `json:"verified"`
	StatusesCount   int         `json:"statuses_count"`
	Lang            interface{} `json:"lang"`
	Status          struct {
		CreatedAt string `json:"created_at"`
		ID        int64  `json:"id"`
		IDStr     string `json:"id_str"`
		Text      string `json:"text"`
		Truncated bool   `json:"truncated"`
		Entities  struct {
			Hashtags     []interface{} `json:"hashtags"`
			Symbols      []interface{} `json:"symbols"`
			UserMentions []interface{} `json:"user_mentions"`
			Urls         []interface{} `json:"urls"`
		} `json:"entities"`
		Source               string      `json:"source"`
		InReplyToStatusID    interface{} `json:"in_reply_to_status_id"`
		InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
		InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
		InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
		InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
		Geo                  interface{} `json:"geo"`
		Coordinates          interface{} `json:"coordinates"`
		Place                interface{} `json:"place"`
		Contributors         interface{} `json:"contributors"`
		IsQuoteStatus        bool        `json:"is_quote_status"`
		RetweetCount         int         `json:"retweet_count"`
		FavoriteCount        int         `json:"favorite_count"`
		Favorited            bool        `json:"favorited"`
		Retweeted            bool        `json:"retweeted"`
		Lang                 string      `json:"lang"`
	} `json:"status"`
	ContributorsEnabled            bool        `json:"contributors_enabled"`
	IsTranslator                   bool        `json:"is_translator"`
	IsTranslationEnabled           bool        `json:"is_translation_enabled"`
	ProfileBackgroundColor         string      `json:"profile_background_color"`
	ProfileBackgroundImageURL      interface{} `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS interface{} `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool        `json:"profile_background_tile"`
	ProfileImageURL                string      `json:"profile_image_url"`
	ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
	ProfileLinkColor               string      `json:"profile_link_color"`
	ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string      `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
	HasExtendedProfile             bool        `json:"has_extended_profile"`
	DefaultProfile                 bool        `json:"default_profile"`
	DefaultProfileImage            bool        `json:"default_profile_image"`
	Following                      interface{} `json:"following"`
	FollowRequestSent              interface{} `json:"follow_request_sent"`
	Notifications                  interface{} `json:"notifications"`
	TranslatorType                 string      `json:"translator_type"`
}
type twitterAPIErrors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// twitter User API リクエストを取ってくる
// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/follow-search-get-users/api-reference/get-users-show
func twitterUserAPIRequst() (*http.Request, error) {
	twitterUserAPIRequest, err := http.NewRequest(http.MethodGet, "https://api.twitter.com/1.1/users/show.json", nil)

	bearerToken := os.Getenv("BEARER_TOKEN")
	twitterUserAPIRequest.Header.Add("authorization", "Bearer "+bearerToken)

	return twitterUserAPIRequest, err
}

// HandleRegistByTwitterName はscreen_name(@hoge)でTwitterIDをSQLに格納する
func HandleRegistByTwitterName(c echo.Context) error {
	req, err := twitterUserAPIRequst()
	if err != nil {
		return c.String(http.StatusBadRequest, "Twitter API Request err")
	}
	screenName := c.Param("screen_name")

	// URLパラメータ作成
	params := req.URL.Query()
	params.Add("screen_name", screenName)
	req.URL.RawQuery = params.Encode()

	// Twitter APIを叩く
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)

	// URLパラメータ削除
	params.Del("screen_name")
	req.URL.RawQuery = params.Encode()

	fmt.Println(req)

	if err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Twitter API Request err")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Cannot read body")
	}

	// jsonにする
	userDataJSON := new(twitterUser)
	if err := json.Unmarshal(body, userDataJSON); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "JSON parse error")
	}

	if userDataJSON.Name == "" {
		return c.String(http.StatusBadRequest, "That user does not exist")
	}
	fmt.Println(userDataJSON)

	// id, screen_name, profile_image_url_https(image_url) をSQLに保存
	insert, err := db.Prepare("INSERT INTO twitter_user(id, screen_name, image_url) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	insert.Exec(userDataJSON.ID, userDataJSON.ScreenName, userDataJSON.ProfileImageURLHTTPS)

	return c.String(http.StatusOK, userDataJSON.Name+" is registration completed")
}
