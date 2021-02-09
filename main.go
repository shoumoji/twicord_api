package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

// type twitterUser struct {
// name            string `json:"name"`
// screenName      string `json:"screen_name"`
// userID          uint   `json:"id"`
// includeEntities bool   `json:"include_entites`
// profileImageURL string `json:"profile_image_url_https"`
// }

type twitterUser struct {
	ID              int64       `json:"id"`
	IDStr           string      `json:"id_str"`
	Name            string      `json:"name"`
	ScreenName      string      `json:"screen_name"`
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

func main() {
	e := echo.New()

	e.POST("/regist/twitter/:username", HandleRegistByTwitterName)
	e.Start(":8000")
}

// Regist twitter ID by twitter screen name
func HandleRegistByTwitterName(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bearerToken := os.Getenv("BEARER_TOKEN")

	twitterUserAPIURL := "https://api.twitter.com/1.1/users/show.json"
	req, err := http.NewRequest("GET", twitterUserAPIURL, nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "URL cannnot make")
	}
	params := req.URL.Query()
	userName := c.Param("username")
	params.Add("screen_name", userName)
	req.URL.RawQuery = params.Encode()
	fmt.Println(req.URL.String())

	req.Header.Add("authorization", "Bearer "+bearerToken)
	if err != nil {
		return c.String(http.StatusBadRequest, "Unknown twitter screen name")
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusBadRequest, "Twitter API Request err")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	userDataJSON := new(twitterUser)
	if err := json.Unmarshal(body, userDataJSON); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "JSON parse error")
	}

	if userDataJSON.Name == "" {
		return c.String(http.StatusBadRequest, "That user does not exist")
	}
	fmt.Println(userDataJSON)

	// TODO: screen_nameと登録日をSQLに保存
	return c.String(http.StatusOK, userDataJSON.Name+" is registration completed")
}
