package dgjtalk

// import "github.com/bwmarrin/dgvoice"
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

const (
	channels     = 2
	samplingRate = 48000
	frameSize    = 960
)

func New(options map[string]string) DgJTalk {
	return DgJTalk{
		options:        options,
		VoiceFile:      "/usr/share/hts-voice/nitech-jp-atr503-m001/nitech_jp_atr503_m001.htsvoice",
		DictionaryPath: "/var/lib/mecab/dic/open-jtalk/naist-jdic",
	}
}

type (
	DgJTalk struct {
		options        map[string]string
		VoiceFile      string
		DictionaryPath string
	}
)

func (d DgJTalk) Talk(v *discordgo.VoiceConnection, text string, stop <-chan struct{}) {

	jTalkArg := []string{}
	// vice file
	jTalkArg = append(jTalkArg, "-m", d.VoiceFile)
	// dictionary path
	jTalkArg = append(jTalkArg, "-x", d.DictionaryPath)
	// sampling rate
	jTalkArg = append(jTalkArg, "-s", strconv.Itoa(samplingRate))
	// Other args
	for k, v := range d.options {
		jTalkArg = append(jTalkArg, fmt.Sprintf("-%s", k), v)
	}
	// pipe command
	jTalkArg = append(jTalkArg, "-ow", "/dev/stdout")

	jtalkCmd := exec.Command("open_jtalk", jTalkArg...)

	ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(samplingRate), "-ac", strconv.Itoa(channels), "pipe:1")

	ffmpegout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// pipe commands
	r, w := io.Pipe()
	jtalkCmd.Stdout = w
	ffmpegCmd.Stdin = r

	in, err := jtalkCmd.StdinPipe()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	io.WriteString(in, text)
	in.Close()

	err = jtalkCmd.Start()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	go func() {
		jtalkCmd.Wait()
		r.Close()
	}()

	err = ffmpegCmd.Start()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	go func() {
		<-stop
		jtalkCmd.Process.Kill()
		ffmpegCmd.Process.Kill()
	}()

	err = v.Speaking(true)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	defer func() {
		err := v.Speaking(false)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		dgvoice.SendPCM(v, send)
		close <- true
	}()

	for {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			return
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return
		}
	}

	return
}
