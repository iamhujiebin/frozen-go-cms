package cron

import "frozen-go-cms/cron/discord_cron"

func Init() {
	discord_cron.GuildFeed()
}
