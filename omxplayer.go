package omxplayer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	"github.com/godbus/dbus"
)

const (
	// Status codes.
	statusStopped  = 0
	statusStarting = 1
	statusPaused   = 2
	statusPlaying  = 3
	statusClosing  = 4
	statusError    = -1

	// DBUS constants
	dbAddressFile = "/tmp/omxplayerdbus."
	dbDestination = "org.mpris.MediaPlayer2.omxplayer"
	dbPath        = "/org/mpris/MediaPlayer2"
	dbInterface   = "org.mpris.MediaPlayer2.Player."
)

// OMXPlayer is a go interface for raspbian's OMXPlayer.
type OMXPlayer struct {
	testing      bool
	status       int
	audioOutput  string
	playbackRate float64
	closing      chan error
	cmd          spyCmd
	dbConn       *dbus.Conn
}

// spyCmd is an Mock/Spy interface for exec.Cmd, for testing purposes.
type spyCmd interface {
	Start() error
	Wait() error
}

// NewOMXPlayer returns a pointer to a new OMXPlayer instance.
// cmd should always be nil, unless a mock is being sent in for a test.
func NewOMXPlayer(testing bool, debug bool, filename string, audioOutput string, cmd spyCmd) (*OMXPlayer, error) {
	if testing && cmd == nil {
		return nil, fmt.Errorf("cmd cannot be nil if testing")
	}

	omx := OMXPlayer{
		testing: testing,
		cmd:     cmd,
	}

	// OMXPlayer command flags.
	args := []string{
		"--blank",             // Set the video background colour to black.
		"--adev", audioOutput, // Set audio output device.
		filename, // File to be played.
	}

	// TODO: Check for loop
	// if loopCondition {
	// 	args = append(args, "--loop")
	// }

	// cmd is only set for testing, otherwise a normal exec.Cmd will be created.
	if !testing {
		cm := exec.Command("omxplayer", args...)
		cm.Stderr = os.Stderr
		if debug {
			cm.Stdout = os.Stdout
		}
		var c spyCmd = cm
		omx.cmd = c

		// Setup dbus connection if not testing
		user, err := user.Current()
		if err != nil {
			// TODO: wrap error.
			return nil, err
		}

		// Get the OMXPlayer session dbus address.
		dbAddr, err := ioutil.ReadFile(dbAddressFile + user.Username)
		if err != nil {
			// TODO: wrap error
			return nil, err
		}

		os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(dbAddr))
		conn, err := dbus.SessionBus()
		if err != nil {
			// TODO: wrap error
			return nil, err
		}

		omx.dbConn = conn
	}

	return &omx, nil
}

// Open starts the OMXPlayer process.
func (o *OMXPlayer) Open(waiting chan string) error {
	o.status = statusStarting

	// Close OMXPlayer if it's already running
	if err := o.Close(); err != nil {
		// TODO wrap error here with the new errors thing
		return err
	}

	if err := o.cmd.Start(); err != nil {
		o.status = statusError
		return err
	}

	o.status = statusPlaying
	o.playbackRate = 1

	// Listen for when OMXPlayer ends in a new goroutine
	go o.wait(waiting)

	return nil
}

func (o *OMXPlayer) wait(waiting chan string) error {
	// Block till the command/process is finished.
	err := o.cmd.Wait()
	prevStatus := o.status
	o.status = statusStopped
	// TODO: remove this condition once dbus has been mocked.
	if !o.testing {
		o.dbConn.Close()
	}
	if prevStatus == statusClosing {
		o.closing <- err
	} else {
		// If OMXPlayer closed naturally, by reaching the end of the file,
		// send a message back to Player to trigger the next item in the playlist.
		waiting <- "next"
	}
	return err
}

// Close ends the OMXPlayer process.
func (o *OMXPlayer) Close() error {
	if o.status != statusPlaying {
		return nil
	}

	o.status = statusClosing
	o.closing = make(chan error)
	defer close(o.closing)

	o.status = statusStopped
	return nil
}

// dbusSend sends dbus messages to the running instance of OMXPlayer
func (o *OMXPlayer) dbusSend(method string, args ...interface{}) (string, error) {
	// Skip this if testing. I'll have to mock the dbus package to test this properly.
	// TODO: mock dbus package to test this properly.
	if o.testing {
		return "", nil
	}

	obj := o.dbConn.Object(dbDestination, dbPath)
	call := obj.Call(dbInterface+method, 0, args...)
	if call.Err != nil {
		// TODO: wrap error
		return "", call.Err
	}

	res := fmt.Sprintf("%v\n", call.Body)
	return string(res), nil
}

// Play the video. If the video is playing, it has no effect, if it is paused it will play from current position.
func (o *OMXPlayer) Play() error {
	if o.status == statusPlaying {
		return nil
	}

	// TODO: test this properly by mocking dbus package.
	if o.testing {
		o.status = statusPlaying
		return nil
	}

	_, err := o.dbusSend("Play", dbus.FlagNoAutoStart)
	if err != nil {
		return err
	}
	o.status = statusPlaying
	return nil
}

// Pause the video. If the video is playing, it will be paused, if it is paused it will stay in pause (no effect).
func (o *OMXPlayer) Pause() error {
	if o.status == statusPaused {
		return nil
	}

	// TODO: test this properly by mocking dbus package.
	if o.testing {
		o.status = statusPaused
		return nil
	}

	_, err := o.dbusSend("Pause", dbus.FlagNoAutoStart)
	if err != nil {
		return err
	}
	o.status = statusPaused
	return nil
}

// Stop stops playback and quits the application.
func (o *OMXPlayer) Stop() {

}

// Chapter seeks the video position to a specific chapter and returns the current chapter.
func (o *OMXPlayer) Chapter(chapter int) int {
	return -1
}

// Seek seeks the video forward or backward in microseconds and returns the current position of the video.
func (o *OMXPlayer) Seek(microseconds int64, whence int) (int64, error) {
	return -1, nil
}

// Position seeks the video to a certain position in the video, in microseconds, and returns the current positon of the video.
func (o *OMXPlayer) Position(microseconds int64) int64 {
	return -1
}

// PlaybackRate sets the playback rate of the video and returns the current playback rate.
func (o *OMXPlayer) PlaybackRate(rate float64) float64 {
	return 0
}

// AudioStream sets the audio stream and returns the current audio stream.
func (o *OMXPlayer) AudioStream(index int) int {
	return -1
}

// VideoStream sets the video stream and returns the current video stream.
func (o *OMXPlayer) VideoStream(index int) int {
	return -1
}

// SubtitleStream sets the subtitle stream and returns the current subtitle stream.
func (o *OMXPlayer) SubtitleStream(index int) int {
	return -1
}

// Volume sets the volume and returns the current volume.
func (o *OMXPlayer) Volume(volume float64) float64 {
	return -1
}
