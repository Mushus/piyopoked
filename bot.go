package piyopoke

import (
	"log"
	"os"
	"sync"

	"github.com/Mushus/piyopoked/chat"
	"github.com/Mushus/piyopoked/dgjtalk"
	"github.com/bwmarrin/discordgo"
)

// Bot ぴよぽけボット
type Bot interface {
	Start() error
	Assign(task interface{})
	SetCallName(callname string)
	Command(command string)
	SetChatClient(client *chat.Client)
}

// NewBot ぴよぽけボットを作成する
func NewBot(callName string, token string) Bot {
	jtalk := dgjtalk.New(map[string]string{
		"a":  "0.45",
		"fm": "-1",
		"jf": "2",
		"r":  "1.1",
	})
	jtalk.VoiceFile = "/usr/share/hts-voice/mei/mei_bashful.htsvoice"
	return &bot{
		exitSignal:   make(chan error),
		jtalk:        jtalk,
		RWMutex:      new(sync.RWMutex),
		state:        offlineState,
		callName:     callName,
		token:        token,
		logger:       log.New(os.Stdout, "[piyopoke] ", log.LstdFlags|log.Lshortfile),
		vcStopSignal: make(chan struct{}),
	}
}

// ボット
type bot struct {
	exitSignal chan error
	// ボット変更時のロック
	*sync.RWMutex
	// ステータス
	state botState
	// ボットの名前
	callName string
	// トークン
	token string
	// discord
	discord *discordgo.Session
	// VCの接続
	voiceConn *discordgo.VoiceConnection
	// VCストップ
	vcStopSignal chan struct{}
	// ロガー
	logger *log.Logger
	// jtalk
	jtalk dgjtalk.DgJTalk
	// チャットクライアント
	chat *chat.Client
}

func (b *bot) Start() error {
	// discord の立ち上げ
	err := func() error {
		b.Lock()
		defer b.Unlock()
		var err error
		b.discord, err = discordgo.New()
		if err != nil {
			return err
		}

		b.discord.Token = b.token
		b.discord.AddHandler(b.onMessageCreate)
		return nil
	}()
	if err != nil {
		return err
	}

	b.Assign(WakeUpTask{})
	return <-b.exitSignal
}

func (b *bot) SetChatClient(client *chat.Client) {
	b.chat = client
}

func (b *bot) Assign(task interface{}) {
	b.Lock()
	b.logger.Printf("task: %#v", task)
	b.state = b.state(b, task)
	b.logger.Printf("new state: %T", b.state)
	defer b.Unlock()
}

func (b *bot) SetCallName(callname string) {
	b.Lock()
	b.callName = callname
	b.Unlock()
}
