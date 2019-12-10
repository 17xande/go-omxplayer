package omxplayer

import (
	"io"
	"testing"
)

// mockCmd is a mock of exec.Cmd.
type mockCmd struct {
	status string
	Stdout io.Writer
}

func (c *mockCmd) Start() error {
	c.status = "started"
	return nil
}

func (c *mockCmd) Wait() error {
	c.status = "ended"
	return nil
}

func setupOMX() (*OMXPlayer, error) {
	var cmd spyCmd = &mockCmd{}
	return NewOMXPlayer(true, true, "video.mp4", "hdmi", cmd)
}

func TestNewOMXPlayer(t *testing.T) {
	if _, err := NewOMXPlayer(true, true, "video.mp4", "hdmi", nil); err != nil && err.Error() != "cmd cannot be nil if testing" {
		t.Error("Should have returned correct error message, instead got: ", err)
	}

	omx, err := setupOMX()

	if err != nil {
		t.Error("Could not initialize OMXPlayer", err)
	}

	if !omx.testing {
		t.Error("OMXPlayer should have 'testing' set to true")
	}
}

func TestOpen(t *testing.T) {
	omx, _ := setupOMX()
	waiting := make(chan string)
	defer close(waiting)
	if err := omx.Open(waiting); err != nil {
		t.Error("Could not open OMXPlayer: ", err)
	}
	done := <-waiting
	if done != "next" {
		t.Error("Expected 'next' got: ", done)
	}
}

func TestClose(t *testing.T) {
	omx, _ := setupOMX()

	// Attempt to close while already closed.
	if err := omx.Close(); err != nil {
		t.Error("Expected no error from OMXPlayer.Close() but instead got: ", err)
	}

	// Simulate playing status and then close.
	omx.status = statusPlaying
	_ = omx.Close()
	if omx.status != statusStopped {
		t.Error("Expected OMXPlayer.status to be 0 (statusStopped), instead got: ", omx.status)
	}
}

func TestPlay(t *testing.T) {
	omx, _ := setupOMX()
	omx.status = statusPaused
	if err := omx.Play(); err != nil {
		t.Error("Got error in Play(): ", err)
	}
	if omx.status != statusPlaying {
		t.Error("Expected OMXPlayer.status to be 3, instead got: ", omx.status)
	}
}

func TestPause(t *testing.T) {
	omx, _ := setupOMX()
	omx.status = statusPlaying
	if err := omx.Pause(); err != nil {
		t.Error("Got error in Pause(): ", err)
	}
	if omx.status != statusPaused {
		t.Error("Expected OMXPlayer.status to be 2, instead got: ", omx.status)
	}
}
