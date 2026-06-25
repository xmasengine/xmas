package xgal

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// SampleRate is the fixed sample rate used for all audio.
const SampleRate = 44100

var (
	audioCtx     *audio.Context
	audioCtxOnce sync.Once
)

func audioContext() *audio.Context {
	audioCtxOnce.Do(func() {
		audioCtx = audio.NewContext(SampleRate)
	})
	return audioCtx
}

func decodeAudio(name string, reader io.Reader) (io.ReadSeeker, error) {
	switch ext := strings.ToLower(filepath.Ext(name)); ext {
	case ".wav":
		return wav.DecodeF32(reader)
	case ".mp3":
		return mp3.DecodeF32(reader)
	case ".ogg":
		return vorbis.DecodeF32(reader)
	default:
		return nil, fmt.Errorf("xgal: unsupported audio format: %s", name)
	}
}

// Clip is a short sound effect loaded with [Sample].
// It plays once. Methods: [Clip.Play], [Clip.Stop], [Clip.IsPlaying], [Clip.Volume].
type Clip struct {
	player *audio.Player
}

// Sample loads an audio file from fsys as a [Clip].
// Supported formats: WAV (.wav), MP3 (.mp3), OGG Vorbis (.ogg).
func Sample(fsys fs.FS, name string) (*Clip, error) {
	buf, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, err
	}

	ctx := audioContext()
	stream, err := decodeAudio(name, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	player, err := ctx.NewPlayerF32(stream)
	if err != nil {
		return nil, err
	}

	return &Clip{player: player}, nil
}

// Play starts playback from the beginning.
func (c *Clip) Play() {
	c.player.Rewind()
	c.player.Play()
}

// Stop pauses playback.
func (c *Clip) Stop() {
	c.player.Pause()
}

// IsPlaying reports whether the clip is currently playing.
func (c *Clip) IsPlaying() bool {
	return c.player.IsPlaying()
}

// Volume sets the playback volume (0.0 silences, 1.0 is full).
func (c *Clip) Volume(v float64) {
	c.player.SetVolume(v)
}

// Song is a looping music track loaded with [Track].
// It loops forever. Methods: [Song.Play], [Song.Stop], [Song.IsPlaying], [Song.Volume].
type Song struct {
	player *audio.Player
}

// Track loads an audio file from fsys as a looping [Song].
// Supported formats: WAV (.wav), MP3 (.mp3), OGG Vorbis (.ogg).
func Track(fsys fs.FS, name string) (*Song, error) {
	buf, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, err
	}

	ctx := audioContext()
	stream, err := decodeAudio(name, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	length, err := stream.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if _, err := stream.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	player, err := ctx.NewPlayerF32(audio.NewInfiniteLoopF32(stream, length))
	if err != nil {
		return nil, err
	}

	return &Song{player: player}, nil
}

// Play starts playback from the beginning.
func (s *Song) Play() {
	s.player.Rewind()
	s.player.Play()
}

// Stop pauses playback.
func (s *Song) Stop() {
	s.player.Pause()
}

// IsPlaying reports whether the song is currently playing.
func (s *Song) IsPlaying() bool {
	return s.player.IsPlaying()
}

// Volume sets the playback volume (0.0 silences, 1.0 is full).
func (s *Song) Volume(v float64) {
	s.player.SetVolume(v)
}
