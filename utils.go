package gosdbot

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

//Check prints the results of an error if it exists
func Check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var myClient = &http.Client{Timeout: 10 * time.Second} //to help handle getting json from the web
//GetJSON ...
func GetJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

//RunIt takes in a chance (1:chance) and rolls to see if the number generator for that range hits 0
func RunIt(chance int) bool {
	rnged, _ := rand.Int(rand.Reader, big.NewInt(int64(chance)))
	if rnged.Int64() == int64(0) {
		return true
	}
	return false
}

//GetServer retrieve the guild/server by name from the current message and session
func GetServer(session *discordgo.Session, msg *discordgo.MessageCreate) string {
	guild := GetGuild(session, msg)
	if guild != nil {
		return guild.Name
	}
	return ""
}

//GetGuild retrieve the actual server/guild object
func GetGuild(session *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.Guild {
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		channel, err = session.Channel(msg.ChannelID)
		if err != nil {
			Check(err)
			return nil
		}
	}

	// Attempt to get the guild from the state,
	// If there is an error, fall back to the restapi.
	guild, err := session.State.Guild(channel.GuildID)
	if err != nil {
		guild, err = session.Guild(channel.GuildID)
		if err != nil {
			Check(nil)
			return nil
		}
	}
	return guild
}

//GetRoles retrieves the roles of a user based on the message the sent+the current session
func GetRoles(session *discordgo.Session, msg *discordgo.MessageCreate) []string {
	currentGuild := GetGuild(session, msg)

	user, err := session.GuildMember(currentGuild.ID, msg.Author.ID)
	Check(err)
	userRoles := user.Roles
	var roleList []string
	for _, role := range userRoles {
		roleObject, err := session.State.Role(currentGuild.ID, role)
		Check(err)
		roleList = append(roleList, roleObject.Name)
	}
	return roleList
}

//In checks if s is in the list
func In(s string, list []string) bool {
	for _, item := range list {
		if s == item {
			return true
		}
	}
	return false
}

//GetEmoji retrieves server emoji by name if it exists.
func GetEmoji(session *discordgo.Session, msg *discordgo.MessageCreate, name string) *discordgo.Emoji {
	guild := GetGuild(session, msg)
	guildEmojis := guild.Emojis
	if guildEmojis != nil {
		for _, guildEmoji := range guildEmojis {
			if guildEmoji.Name == name {
				return guildEmoji
			}
		}
	}

	return nil
}

//CheckWebHooks parses through a list of webhooks given the channel to see if the target is present
func CheckWebHooks(session *discordgo.Session, message *discordgo.MessageCreate, target string) bool {
	webhooks, err := session.ChannelWebhooks(message.ChannelID)
	Check(err)

	for _, hook := range webhooks {
		if target == hook.Name {
			return true
		}
	}
	return false
}
