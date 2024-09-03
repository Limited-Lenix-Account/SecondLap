package discord

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	disgo "github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

const (
	avatarURL  = "" // put ur avi url
	WebhookURL = "" // put ur discord url
)

type DiscordData struct {
	Email      string
	EntryCount int
}

// thank you frankwhite
func (d *DiscordData) SendEmbed() error {

	snowflakeReg := regexp.MustCompile("[0-9]{18,19}")
	it, _ := strconv.Atoi(snowflakeReg.FindString(WebhookURL))
	split := strings.Split(WebhookURL, "/")

	var (
		sf      = snowflake.ID(it)
		channel = split[len(split)-1]
	)

	client := disgo.New(sf, channel)
	defer client.Close(context.TODO())

	discordEmbed := discord.NewEmbedBuilder()

	discordEmbed.SetTitle("Entry Submitted!")
	// discordEmbed.SetURL(productUrl)
	// discordEmbed.SetDescription("Click the link above to view inventory page")
	discordEmbed.SetThumbnail("https://pbs.twimg.com/profile_images/917516130926739456/qn32fLJ9_400x400.jpg")
	discordEmbed.SetFooterText("Made By Lenix")
	discordEmbed.SetFooterIcon(avatarURL)
	discordEmbed.SetTimestamp(time.Now().Local())
	//
	discordEmbed.AddField("Site", "KTM", true)
	discordEmbed.AddField("Email 📧", d.Email, true)
	discordEmbed.AddField("Entry Count 📌", fmt.Sprintf("%d", d.EntryCount), true)

	_, err := client.Rest().CreateWebhookMessage(sf, channel, discord.WebhookMessageCreate{
		Username:  "Lenix's Helper",
		AvatarURL: avatarURL,
		Embeds:    []discord.Embed{discordEmbed.Build()},
	}, true, snowflake.ID(0))

	if err != nil {
		log.Printf("Error Sending Discord Embed: %s\n", err)
		return err
	}
	return nil
}
