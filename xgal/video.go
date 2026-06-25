package xgal

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"sync"
	"time"

	"github.com/gen2brain/mpeg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Video is a playable MPEG video stream loaded with [Stream].
// Methods: [Video.Draw], [Video.Play], [Video.Stop], [Video.IsPlaying],
// [Video.Frame], [Video.Bounds], [Video.HasEnded], [Video.Close].
type Video struct {
	mpg   *mpeg.MPEG
	frame *Surface

	audioPlayer *audio.Player

	src       io.ReadCloser
	refTime   time.Time
	closeOnce sync.Once
	m         sync.Mutex
}

// Stream loads an MPEG video file from fsys as a [Video].
func Stream(fsys fs.FS, name string) (*Video, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}

	mpg, err := mpeg.New(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	if mpg.NumVideoStreams() == 0 {
		f.Close()
		return nil, fmt.Errorf("xgal: no video streams in %s", name)
	}
	if !mpg.HasHeaders() {
		f.Close()
		return nil, fmt.Errorf("xgal: missing headers in %s", name)
	}

	v := &Video{
		mpg:   mpg,
		frame: NewSurface(mpg.Width(), mpg.Height()),
		src:   f,
	}

	if mpg.NumAudioStreams() > 0 {
		ctx := audioContext()
		if ctx.SampleRate() != mpg.Samplerate() {
			f.Close()
			return nil, fmt.Errorf("xgal: video audio %d Hz != %d Hz", mpg.Samplerate(), ctx.SampleRate())
		}
		mpg.SetAudioFormat(mpeg.AudioF32N)
		ap, err := ctx.NewPlayerF32(&mpegAudio{audio: mpg.Audio(), m: &v.m})
		if err != nil {
			f.Close()
			return nil, err
		}
		v.audioPlayer = ap
	}

	return v, nil
}

// Draw updates to the current frame and draws it onto screen, fitting it
// within the screen bounds while maintaining the aspect ratio.
func (v *Video) Draw(screen *Surface) {
	v.m.Lock()
	defer v.m.Unlock()

	pos := v.playbackPos()
	video := v.mpg.Video()
	if video.HasEnded() {
		return
	}

	d := 1 / v.mpg.Framerate()
	var mpegFrame *mpeg.Frame
	for video.Time()+d <= pos && !video.HasEnded() {
		mpegFrame = video.Decode()
	}
	if mpegFrame == nil {
		return
	}

	rgba := mpegFrame.RGBA()
	v.frame.WritePixels(rgba.Pix)

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	fw, fh := v.frame.Bounds().Dx(), v.frame.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	s := math.Min(float64(sw)/float64(fw), float64(sh)/float64(fh))
	op.GeoM.Scale(s, s)
	op.GeoM.Translate((float64(sw)-float64(fw)*s)/2, (float64(sh)-float64(fh)*s)/2)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(v.frame, op)
}

// Frame returns the current decoded frame as a [Surface].
func (v *Video) Frame() *Surface {
	return v.frame
}

// Play starts playback.
func (v *Video) Play() {
	v.m.Lock()
	defer v.m.Unlock()

	if v.mpg.HasEnded() {
		v.mpg.Rewind()
	}

	if v.audioPlayer != nil {
		if v.audioPlayer.IsPlaying() {
			return
		}
		go v.audioPlayer.Play()
		return
	}

	if v.refTime != (time.Time{}) {
		return
	}
	v.refTime = time.Now()
}

// Stop pauses playback.
func (v *Video) Stop() {
	if v.audioPlayer != nil {
		v.audioPlayer.Pause()
		return
	}
	v.refTime = time.Time{}
}

// IsPlaying reports whether the video is currently playing.
func (v *Video) IsPlaying() bool {
	if v.audioPlayer != nil {
		return v.audioPlayer.IsPlaying()
	}
	return v.refTime != (time.Time{})
}

func (v *Video) playbackPos() float64 {
	if v.audioPlayer != nil {
		return v.audioPlayer.Position().Seconds()
	}
	if v.refTime != (time.Time{}) {
		return time.Since(v.refTime).Seconds()
	}
	return 0
}

// Close closes the video source file.
func (v *Video) Close() error {
	var err error
	v.closeOnce.Do(func() {
		if v.audioPlayer != nil {
			v.audioPlayer.Pause()
		}
		err = v.src.Close()
	})
	return err
}

// Bounds returns the video dimensions.
func (v *Video) Bounds() (int, int) {
	return v.mpg.Width(), v.mpg.Height()
}

// HasEnded reports whether the video has finished playing.
func (v *Video) HasEnded() bool {
	return v.mpg.HasEnded()
}

// mpegAudio implements io.Reader for the audio player by decoding MPEG audio
// samples. The shared mutex m prevents concurrent access to the MPEG decoder.
type mpegAudio struct {
	audio     *mpeg.Audio
	leftovers []byte
	m         *sync.Mutex
}

func (a *mpegAudio) Read(buf []byte) (int, error) {
	a.m.Lock()
	defer a.m.Unlock()

	var readBytes int
	if len(a.leftovers) > 0 {
		n := copy(buf, a.leftovers)
		readBytes += n
		buf = buf[n:]
		copy(a.leftovers, a.leftovers[n:])
		a.leftovers = a.leftovers[:len(a.leftovers)-n]
	}

	for len(buf) > 0 && !a.audio.HasEnded() {
		samples := a.audio.Decode()
		if samples == nil {
			break
		}

		bs := samples.Bytes()
		n := copy(buf, bs)
		readBytes += n
		buf = buf[n:]

		if n < len(bs) {
			a.leftovers = append(a.leftovers, bs[n:]...)
			break
		}
	}

	if a.audio.HasEnded() {
		return readBytes, io.EOF
	}
	return readBytes, nil
}
