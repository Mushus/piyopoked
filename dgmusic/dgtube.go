package dgmusic

// import "github.com/bwmarrin/dgvoice"
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"os"
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

func NewDgTube() DgTube {
	return DgTube{}
}

type (
	DgTube struct {
	}
)

func (d DgTube) Play(v *discordgo.VoiceConnection, url url.URL, stop <-chan struct{}) {

	youtubeArg := []string{url.String(), "-o", "-"}

	youtubeCmd := exec.Command("youtube-dl", youtubeArg...)

	ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-af", "volume=-10dB", "-ar", strconv.Itoa(samplingRate), "-ac", strconv.Itoa(channels), "pipe:1")

	ffmpegout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// pipe commands
	r, w := io.Pipe()
	youtubeCmd.Stderr = os.Stderr
	youtubeCmd.Stdout = w
	ffmpegCmd.Stdin = r

	err = youtubeCmd.Start()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	go func() {
		youtubeCmd.Wait()
		r.Close()
	}()

	err = ffmpegCmd.Start()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	go func() {
		<-stop
		youtubeCmd.Process.Kill()
		ffmpegCmd.Process.Kill()
	}()

	err = v.Speaking(true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	defer func() {
		err := v.Speaking(false)
		if err != nil {
			fmt.Printf("%v\n", err)
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
			fmt.Printf("%v\n", err)
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
