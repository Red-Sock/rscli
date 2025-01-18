package animation

import (
	"bytes"
	"os"
	"time"

	"go.redsock.ru/rerrors"
)

var (
	ErrNoFrames          = rerrors.New("no frames")
	ErrFrameSizeNotMatch = rerrors.New("frame size not match")
)

type opt func(a *Animation) error

func WithWriter(in *os.File, out *os.File) opt {
	return func(a *Animation) error {
		a.in = in
		a.out = out
		return nil
	}
}

func WithDuration(dur time.Duration) opt {
	return func(a *Animation) error {
		a.duration = dur
		return nil
	}
}

func WithStrFrames(framesStr ...string) opt {
	frames := make([][]byte, 0, len(framesStr))

	for _, f := range framesStr {
		frames = append(frames, []byte(f))
	}

	return WithFrames(frames...)
}

func WithFrames(frames ...[]byte) opt {
	return func(a *Animation) error {
		if len(frames) == 0 {
			return ErrNoFrames
		}

		a.width, a.height = bytes.IndexByte(frames[0], '\n'), bytes.Count(frames[0], []byte{'\n'})

		for _, frame := range frames[1:] {
			width, height := bytes.IndexByte(frame, '\n'), bytes.Count(frame, []byte{'\n'})
			if width != a.width || height != a.height {
				return ErrFrameSizeNotMatch
			}
		}
		a.frames = frames

		return nil
	}
}
