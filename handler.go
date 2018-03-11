package piyopoke

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func (b *bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.RLock()
	callName := b.callName
	b.RUnlock()

	content := m.Content
	b.logger.Print(content)
	commandRegexp := regexp.MustCompile(fmt.Sprintf("^%s[\\s　]+(.*)[\\s　]*", regexp.QuoteMeta(callName)))
	res := commandRegexp.FindStringSubmatch(content)
	// コマンドじゃない
	if res == nil {
		return
	}

	command := res[1]
	userID := m.Author.ID
	channelID := m.ChannelID

	b.command(command, userID, channelID)
}
