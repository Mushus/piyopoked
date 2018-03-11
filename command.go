package piyopoke

import (
	"net/url"
	"regexp"
)

var (
	// ExitCmd 終了コマンド
	ExitCmd = regexp.MustCompile("^!exit$")
	// ComeCmd 呼ぶコマンド
	ComeCmd = regexp.MustCompile("^((こっち|コッチ|ここ|ココ)[へ|に]?)?(おいで|[来き](てね?|なさい)|[こ来]い)[！？\\!\\?\\.,。、…]*$")
	// ByeByeCmd バイバイコマンド
	ByeByeCmd = regexp.MustCompile("^(ばいばい|バイバイ|出て|でて|切って|きって)[！？\\!\\?\\.,。、…]*$")
	// JoinCmd 参加コマンド
	JoinMyGuildVCCmd = regexp.MustCompile("^(.+)[にへ]((入|はい)(りなさい|って(くださいね?|ね)?|れ)|(参加|さんか)(し(なさい|て(くださいね?|ね)?|ろ))|おいで|[来き](てね?|なさい)|[こ来]い)[！？\\!\\?\\.,。、…]*$")
	JoinVCCmd        = regexp.MustCompile("^(.+)\\s?/\\s?(.+)[に|へ]((入|はい)(りなさい|って(くださいね?|ね)?|れ)|(参加|さんか)(し(なさい|て(くださいね?|ね)?|ろ))|おいで|[来き](てね?|なさい)|[こ来]い)[！？\\!\\?\\.,。、…]*$")
	TalkCmd          = regexp.MustCompile("^(.+)(って|と)((しゃべ|喋)(れ|ろ|りなさい|って(ください|ね)?)|(話|はな)(して[よね]?|せや?)|[言い](え|って(よね?|ください)?))[！？\\!\\?\\.,。、…]*$")
	StopCmd          = regexp.MustCompile("^([やと止]め(てね?|ろ|てくださいね?)?|ストップ|(黙|だま)(れよ?|って(よ(ね?)?))|(終|おわ)(って(よね?|ください)?|り|れ)|[Ss][Tt][Oo][Pp]|停止|ていし|(五月蝿|煩|うるさ)い)[！？\\!\\?\\.,。、…]*$")
	ExtarnalMusicCmd = regexp.MustCompile("^https://(m|www)\\.youtube\\.com/watch.+|http://www\\.nicovideo\\.jp/watch/[a-z0-9]*$")
	InternalMusicCmd = regexp.MustCompile("^([a-zA-Z0-9_-]+)を((再生|さいせい)して(ね|よ|くれ)?|(流|なが)(せ|して)(よ|くれ)?)$")
)

func (b *bot) Command(command string) {
	b.command(command, "", "")
}

func (b *bot) command(command string, userID string, channelID string) {
	// コマンドからタスクを振る
	if ExitCmd.MatchString(command) {
		b.logger.Print("exit")
		b.Assign(ExitTask{})
		return
	}
	if ComeCmd.MatchString(command) {
		b.Assign(ComeVCTask{
			CmdUserID:    userID,
			CmdChannelID: channelID,
		})
		return
	}
	if ByeByeCmd.MatchString(command) {
		b.Assign(BayBayVCTask{})
		return
	}
	if match := JoinVCCmd.FindStringSubmatch(command); len(match) > 0 {
		b.Assign(JoinVCTask{
			CmdUserID:    userID,
			CmdChannelID: channelID,
			GuildName:    match[1],
			ChannelName:  match[2],
		})
		return
	}
	if match := JoinMyGuildVCCmd.FindStringSubmatch(command); len(match) > 0 {
		b.logger.Println(match)
		b.Assign(JoinVCTask{
			CmdUserID:    userID,
			CmdChannelID: channelID,
			ChannelName:  match[1],
		})
		return
	}
	if match := TalkCmd.FindStringSubmatch(command); len(match) > 0 {
		b.Assign(TalkTask{
			Message: match[1],
		})
		return
	}
	if StopCmd.MatchString(command) {
		b.Assign(StopTask{})
		return
	}
	if match := ExtarnalMusicCmd.FindStringSubmatch(command); len(match) > 0 {
		uri, err := url.Parse(match[0])
		if err != nil {
			return
		}
		b.Assign(ExtarnalMusicTask{
			URL: *uri,
		})
		return
	}
	if match := InternalMusicCmd.FindStringSubmatch(command); len(match) > 0 {
		b.Assign(InternalMusicTask{
			Name: match[1],
		})
		return
	}
	b.Assign(ChatTask{
		ChannelID: channelID,
		UserID:    userID,
		Message:   command,
	})
}
