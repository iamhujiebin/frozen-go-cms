package discord_cron

import (
	"encoding/json"
	"frozen-go-cms/domain/model/discord_m"
	"frozen-go-cms/hilo-common/domain"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"time"
)

func GuildFeed() {
	model := domain.CreateModelNil()
	res := discordFeeds(model) // init run
	if res != nil {
		for _, v := range res.Results.Items {
			for _, v2 := range v.Message.Attachments {
				if err := discord_m.AddDiscordFeed(model, discord_m.DiscordFeed{
					MsgId:       v.Message.ID,
					MsgContent:  v.Message.Content,
					Author:      v.Message.Author.Username,
					Avatar:      v.Message.Author.Avatar,
					AttachId:    v2.ID,
					Url:         v2.URL,
					ProxyUrl:    v2.ProxyURL,
					Width:       v2.Width,
					Height:      v2.Height,
					Size:        v2.Size,
					ContentType: v2.ContentType,
				}); err != nil {
					model.Log.Errorf("AddDiscordFeed fail:%v", err)
				}
			}
		}
	}
	c := cron.New()
	// 1小时刷新一下api
	spec := "0 0 * * * ?"
	_ = c.AddFunc(spec, func() {
		discordFeeds(model)
	})
	c.Start()
}

type DiscordResponse struct {
	Results struct {
		Items []struct {
			Message struct {
				ID        string `json:"id"`
				Type      int    `json:"type"`
				Content   string `json:"content"`
				ChannelID string `json:"channel_id"`
				Author    struct {
					ID               string      `json:"id"`
					Username         string      `json:"username"`
					GlobalName       string      `json:"global_name"`
					Avatar           string      `json:"avatar"`
					Discriminator    string      `json:"discriminator"`
					PublicFlags      int         `json:"public_flags"`
					AvatarDecoration interface{} `json:"avatar_decoration"`
				} `json:"author"`
				Attachments []struct {
					ID          string `json:"id"`
					Filename    string `json:"filename"`
					Size        int    `json:"size"`
					URL         string `json:"url"`
					ProxyURL    string `json:"proxy_url"`
					Width       int    `json:"width"`
					Height      int    `json:"height"`
					ContentType string `json:"content_type"`
				} `json:"attachments"`
				Embeds          []interface{} `json:"embeds"`
				Mentions        []interface{} `json:"mentions"`
				MentionRoles    []interface{} `json:"mention_roles"`
				Pinned          bool          `json:"pinned"`
				MentionEveryone bool          `json:"mention_everyone"`
				Tts             bool          `json:"tts"`
				Timestamp       time.Time     `json:"timestamp"`
				EditedTimestamp interface{}   `json:"edited_timestamp"`
				Flags           int           `json:"flags"`
				Components      []interface{} `json:"components"`
				Reactions       []struct {
					Emoji struct {
						ID   interface{} `json:"id"`
						Name string      `json:"name"`
					} `json:"emoji"`
					Count        int `json:"count"`
					CountDetails struct {
						Burst  int `json:"burst"`
						Normal int `json:"normal"`
					} `json:"count_details"`
					BurstColors []interface{} `json:"burst_colors"`
					MeBurst     bool          `json:"me_burst"`
					Me          bool          `json:"me"`
					BurstCount  int           `json:"burst_count"`
				} `json:"reactions"`
			} `json:"message"`
			ReferenceMessages []interface{} `json:"reference_messages"`
			Type              string        `json:"type"`
			ID                string        `json:"id"`
			Seen              bool          `json:"seen"`
		} `json:"items"`
	} `json:"results"`
	LoadID string `json:"load_id"`
}

func discordFeeds(model *domain.Model) *DiscordResponse {
	url := "https://discord.com/api/v9/guilds/662267976984297473/guild-feed?limit=10&offset=0"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		model.Log.Errorf("discordFeeds fail:%v", err)
		return nil
	}
	req.Header.Add("authorization", "MTEyODYxNDQ4ODQyMjIyMzg4Mg.G22JEY.H324-ZZaqxFpvNLqOOFNc2EHOnnljfcm6eP8m8")

	res, err := client.Do(req)
	if err != nil {
		model.Log.Errorf("discordFeeds fail:%v", err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		model.Log.Errorf("discordFeeds fail:%v", err)
		return nil
	}
	response := new(DiscordResponse)
	err = json.Unmarshal(body, res)
	return response
}
