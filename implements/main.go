package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/Mushus/piyopoked"
	"github.com/Mushus/piyopoked/chat"
	"github.com/carlescere/scheduler"
)

var BGMList = []string{
	"nowornever",
	"nowornever2",
	"noworneverfest",
	"noworneverfest2",
}

type Config struct {
	BotName            string `json:"bot_name"`
	DiscordClientToken string `json:"discord_client_token"`
	ChatAPIKey         string `json:"chat_api_key"`
}

func main() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	chat := chat.NewClient(cfg.ChatAPIKey)

	stop := make(chan struct{})
	bot := piyopoke.NewBot(cfg.BotName, cfg.DiscordClientToken)
	bot.SetChatClient(chat)
	go func() {
		err := bot.Start()
		if err != nil {
			log.Fatalln(err)
		}
		close(stop)
	}()

	rand.Seed(time.Now().UnixNano())

	scheduler.Every().Day().At("22:00").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "時間だよ！ワンドロ始めるよー！"})
	})
	scheduler.Every().Day().At("22:30").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "残り30分です"})
	})
	scheduler.Every().Day().At("22:50").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "残り10分です"})
	})
	scheduler.Every().Day().At("22:55").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "残り5分です"})
	})
	scheduler.Every().Day().At("22:59").Run(func() {
		bgm := BGMList[rand.Intn(len(BGMList))]
		bot.Assign(piyopoke.InternalMusicTask{Name: bgm})
	})
	scheduler.Every().Day().At("23:00").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "ワンドロ終了だよ！タグを付けてツイッターに投稿しようね！"})
	})

	scheduler.Every().Day().At("00:00").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "時間だよ！ワンドロ始めるよー！"})
	})
	scheduler.Every().Day().At("00:30").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "残り30分です"})
	})
	scheduler.Every().Day().At("00:50").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "残り10分です"})
	})
	scheduler.Every().Day().At("00:55").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "残り5分です"})
	})
	scheduler.Every().Day().At("00:59").Run(func() {
		bgm := BGMList[rand.Intn(len(BGMList))]
		bot.Assign(piyopoke.InternalMusicTask{Name: bgm})
	})
	scheduler.Every().Day().At("01:00").Run(func() {
		bot.Assign(piyopoke.TalkTask{Message: "ワンドロ終了だよ！タグを付けてツイッターに投稿しようね！"})
	})

	<-stop
}
