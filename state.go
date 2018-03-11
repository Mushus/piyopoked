package piyopoke

// botState ボットのステータス
type botState func(bot *bot, task interface{}) botState

// onlineState オンライン(非VC)
func onlineState(bot *bot, task interface{}) botState {
	switch t := task.(type) {
	case ExitTask:
		if bot.discord.Close() != nil {
			return onlineState
		}
		close(bot.exitSignal)
	case JoinVCTask:
		if bot.joinVCAct(t.CmdUserID, t.CmdChannelID, t.GuildName, t.ChannelName) != nil {
			return onlineState
		}
		bot.sayHelloAct()
		return waitVCState
	case ComeVCTask:
		if bot.comeFromAct(t.CmdChannelID, t.CmdUserID) != nil {
			return onlineState
		}
		bot.sayHelloAct()
		return waitVCState
	case ChatTask:
		bot.ChatAct(t.ChannelID, t.UserID, t.Message)
		return onlineState
	}
	return onlineState
}

// waitVCState ボイスチャットで待機中
func waitVCState(bot *bot, task interface{}) botState {
	switch t := task.(type) {
	case ExitTask:
		bot.sayByeByeAct()
		if bot.voiceConn.Disconnect() != nil {
			return waitVCState
		}
		if bot.discord.Close() != nil {
			return onlineState
		}
		close(bot.exitSignal)
		return offlineState
	case JoinVCTask:
		bot.sayByeByeAct()
		if bot.voiceConn.Disconnect() != nil {
			return waitVCState
		}
		if bot.joinVCAct(t.CmdUserID, t.CmdChannelID, t.GuildName, t.ChannelName) != nil {
			return onlineState
		}
		return waitVCState
	case ComeVCTask:
		bot.sayByeByeAct()
		if bot.voiceConn.Disconnect() != nil {
			return waitVCState
		}
		if bot.comeFromAct(t.CmdChannelID, t.CmdUserID) != nil {
			return onlineState
		}
		bot.sayHelloAct()
		return waitVCState
	case BayBayVCTask:
		bot.sayByeByeAct()
		if bot.voiceConn.Disconnect() != nil {
			return waitVCState
		}
		return onlineState
	case TalkTask:
		bot.stopAct()
		bot.talkAct(t.Message, bot.vcStopSignal)
		return talkingState
	case ExtarnalMusicTask:
		bot.stopAct()
		bot.ExtarnalMusicAct(t.URL, bot.vcStopSignal)
		return playingState
	case InternalMusicTask:
		bot.stopAct()
		bot.IntarnalMusicAct(t.Name, bot.vcStopSignal)
		return playingState
	case ChatTask:
		bot.ChatAct(t.ChannelID, t.UserID, t.Message)
		return waitVCState
	}
	return waitVCState
}

func talkingState(bot *bot, task interface{}) botState {
	switch t := task.(type) {
	case ExitTask:
		bot.stopAct()
		bot.sayByeByeAct()
		if bot.voiceConn.Disconnect() != nil {
			return waitVCState
		}
		if bot.discord.Close() != nil {
			return onlineState
		}
		close(bot.exitSignal)
		return offlineState
	case ComeVCTask:
		bot.stopAct()
		bot.sayByeByeAct()
		if bot.comeFromAct(t.CmdChannelID, t.CmdUserID) != nil {
			return onlineState
		}
		bot.sayHelloAct()
		return waitVCState
	case StopTask:
		bot.stopAct()
		return waitVCState
	case TalkTask:
		bot.stopAct()
		bot.talkAct(t.Message, bot.vcStopSignal)
		return talkingState
	case ExtarnalMusicTask:
		bot.stopAct()
		bot.ExtarnalMusicAct(t.URL, bot.vcStopSignal)
		return playingState
	case InternalMusicTask:
		bot.stopAct()
		bot.IntarnalMusicAct(t.Name, bot.vcStopSignal)
		return playingState
	case ChatTask:
		bot.ChatAct(t.ChannelID, t.UserID, t.Message)
		return talkingState
	}
	return talkingState
}

func playingState(bot *bot, task interface{}) botState {
	switch t := task.(type) {
	case ExitTask:
		bot.stopAct()
		bot.sayByeByeAct()
		if bot.voiceConn.Disconnect() != nil {
			return waitVCState
		}
		if bot.discord.Close() != nil {
			return onlineState
		}
		close(bot.exitSignal)
		return offlineState
	case ComeVCTask:
		bot.stopAct()
		bot.sayByeByeAct()
		if bot.comeFromAct(t.CmdChannelID, t.CmdUserID) != nil {
			return onlineState
		}
		bot.sayHelloAct()
		return waitVCState
	case StopTask:
		bot.stopAct()
		return waitVCState
	case TalkTask:
		bot.stopAct()
		bot.talkAct(t.Message, bot.vcStopSignal)
		return talkingState
	case ExtarnalMusicTask:
		bot.stopAct()
		bot.ExtarnalMusicAct(t.URL, bot.vcStopSignal)
		return playingState
	case InternalMusicTask:
		bot.stopAct()
		bot.IntarnalMusicAct(t.Name, bot.vcStopSignal)
		return playingState
	case ChatTask:
		bot.ChatAct(t.ChannelID, t.UserID, t.Message)
		return playingState
	}
	return playingState
}

func offlineState(bot *bot, task interface{}) botState {
	switch task.(type) {
	case WakeUpTask:
		bot.discord.Open()
		return onlineState
	}
	return offlineState
}
