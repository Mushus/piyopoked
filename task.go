package piyopoke

import "net/url"

type (
	// WakeUpTask オンライン状態にする
	WakeUpTask struct {
	}
	// JoinVCTask ボイスチャットに入る
	JoinVCTask struct {
		CmdUserID    string
		CmdChannelID string
		GuildName    string
		ChannelName  string
	}
	// ComeVCTask ボイスチャットに来る
	ComeVCTask struct {
		CmdChannelID string
		CmdUserID    string
	}
	// BayBayVCTask ボイスチャットから抜ける
	BayBayVCTask struct {
	}

	// TalkTask 話す
	TalkTask struct {
		Message string
	}

	// ExtarnalMusicTask 外部の音楽を再生する
	ExtarnalMusicTask struct {
		URL url.URL
	}

	// InternalMusicTask 内部の音楽を再生する
	InternalMusicTask struct {
		Name string
	}

	// StopTask 止める
	StopTask struct {
	}

	// ExitTask 終了
	ExitTask struct {
	}

	// ChatTask チャットする
	ChatTask struct {
		ChannelID string
		UserID    string
		Message   string
	}
)
