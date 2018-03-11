package dgmusic

// import "github.com/bwmarrin/dgvoice"
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func NewDgMusic() DgMusic {
	return DgMusic{}
}

type (
	DgMusic struct {
	}
)

func (d DgMusic) Play(v *discordgo.VoiceConnection, filename string, stop <-chan struct{}) {
	ffmpegCmd := exec.Command("ffmpeg", "-i", filepath.Join(".", "voice", filename+".wav"), "-f", "s16le", "-af", "volume=-10dB", "-ar", strconv.Itoa(samplingRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	ffmpegCmd.Stderr = os.Stderr

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	err = ffmpegCmd.Start()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	go func() {
		<-stop
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
