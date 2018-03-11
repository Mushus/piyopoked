package piyopoke

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Mushus/piyopoked/dgmusic"
	"github.com/bwmarrin/discordgo"
)

type feedbackError struct {
	err string
}

func (e feedbackError) FeedbackError() string {
	return e.err
}

func (e feedbackError) Error() string {
	return e.err
}

var (
	// ErrorUserIsNotInVC ユーザーがVCに入っていない
	ErrorUserIsNotInVC = feedbackError{err: "ボイスチャットに入ってないよ？"}
	// ErrorNoSuchGuild ギルドがない
	ErrorNoSuchGuild = feedbackError{err: "そんなギルド見つからないよ？"}
	// ErrorNoSuchChannel ギルドがない
	ErrorNoSuchChannel = feedbackError{err: "そんなチャンネル見つからないよ？"}
	// ErrorGuildNotSelected ギルド名指定なし
	ErrorGuildNotSelected = feedbackError{err: "ギルドが指定されてないよ？"}
)

func feedback(s *discordgo.Session, userID string, channelID string, fe feedbackError) error {
	_, err := s.ChannelMessageSend(channelID, fmt.Sprintf("<@%s> %s", userID, fe.FeedbackError()))
	return err
}

func (b *bot) comeFromAct(cmdChannelID string, userID string) error {
	d := b.discord
	guildID, channelID, err := findUserVCChannelID(d, userID)
	if err != nil {
		if fe, ok := err.(feedbackError); ok {
			feedback(d, userID, cmdChannelID, fe)
		} else {
			b.logger.Println(err)
		}
		return err
	}
	b.voiceConn, err = d.ChannelVoiceJoin(guildID, channelID, false, false)
	return err
}

func (b *bot) joinVCAct(cmdUserID string, cmdChannelID string, guildName string, channelName string) error {
	d := b.discord
	guildID, channelID, err := findSelectedVC(d, cmdUserID, cmdChannelID, guildName, channelName)
	if err != nil {
		if fe, ok := err.(feedbackError); ok {
			feedback(d, cmdUserID, cmdChannelID, fe)
		} else {
			b.logger.Println(err)
		}
		return err
	}
	b.voiceConn, err = d.ChannelVoiceJoin(guildID, channelID, false, false)
	return err
}

func (b *bot) talkAct(message string, stop chan struct{}) error {
	go func() {
		b.jtalk.Talk(b.voiceConn, message, stop)
	}()
	return nil
}

func (b *bot) ExtarnalMusicAct(uri url.URL, stop chan struct{}) error {
	go func() {
		tube := dgmusic.NewDgTube()
		tube.Play(b.voiceConn, uri, stop)
	}()
	return nil
}

func (b *bot) IntarnalMusicAct(name string, stop chan struct{}) error {
	go func() {
		tube := dgmusic.NewDgMusic()
		tube.Play(b.voiceConn, name, stop)
	}()
	return nil
}

func (b *bot) stopAct() {
	close(b.vcStopSignal)
	b.vcStopSignal = make(chan struct{})
}

func (b *bot) sayHelloAct() error {
	b.jtalk.Talk(b.voiceConn, "こんにちはー！", make(chan struct{}))
	return nil
}

func (b *bot) sayByeByeAct() error {
	b.jtalk.Talk(b.voiceConn, "バイバイ！", make(chan struct{}))
	return nil
}

func (b *bot) ChatAct(channelID string, userID string, message string) error {
	d := b.discord
	c := b.chat
	l := b.logger
	go func() {
		user, err := d.User(userID)
		if err != nil {
			l.Println(err)
			return
		}
		time.Sleep(time.Second / 2)
		d.ChannelTyping(channelID)
		resp, err := c.Talk(user.Username, message)
		if err != nil {
			l.Println(err)
			return
		}
		time.Sleep(time.Second / 2)
		_, err = d.ChannelMessageSend(channelID, fmt.Sprintf("<@%s> %s", userID, resp))
	}()
	return nil
}

func findUserVCChannelID(s *discordgo.Session, userID string) (string, string, error) {
	userGuilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		return "", "", err
	}
	for _, ug := range userGuilds {
		guild, err := s.Guild(ug.ID)
		if err != nil {
			return "", "", err
		}

		for _, vs := range guild.VoiceStates {
			if vs.UserID == userID {
				return vs.GuildID, vs.ChannelID, nil
			}
		}
	}
	return "", "", ErrorUserIsNotInVC
}

// 指定されたボイスチャンネルを探す
func findSelectedVC(s *discordgo.Session, userID string, channelID string, guildName string, channelName string) (string, string, error) {
	var targetGuild *discordgo.Guild
	var targetChannel *discordgo.Channel
	// ギルド名が指定ない時
	if guildName == "" {
		// 発言したチャンネルの属してるギルドを見てみる
		channel, err := s.Channel(channelID)
		if err != nil {
			return "", "", fmt.Errorf("cannot find channel: ChannelID %v: %v", channelID, err)
		}
		targetGuildID := channel.GuildID
		// ギルドが存在しないチャンネル
		if targetGuildID == "" {
			return "", "", ErrorGuildNotSelected
		}
		targetGuild, err = s.Guild(targetGuildID)
		fmt.Printf("%+v", channel)
		if err != nil {
			return "", "", fmt.Errorf("cannot find guild: GuildID %v: %v", channel.GuildID, err)
		}
	} else {
		// ユーザーの属してるギルドを探す
		userGuilds, err := s.UserGuilds(100, "", "")
		if err != nil {
			return "", "", fmt.Errorf("cannot find users guild: %v", err)
		}
		// その中で名前が一致してるギルドを探す
		for _, ug := range userGuilds {
			guild, err := s.Guild(ug.ID)
			if err != nil {
				return "", "", fmt.Errorf("cannot find guild: %v", err)
			}
			if guild.Name == guildName {
				targetGuild = guild
				break
			}
		}
	}
	if targetGuild == nil {
		return "", "", ErrorNoSuchGuild
	}
	for _, channel := range targetGuild.Channels {
		if channel.Name == channelName {
			targetChannel = channel
			break
		}
	}
	if targetChannel == nil {
		return "", "", ErrorNoSuchChannel
	}
	return targetGuild.ID, targetChannel.ID, nil
}
