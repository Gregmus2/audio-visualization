package internal

import (
	"github.com/Gregmus2/simple-engine/scenes"
	"github.com/sirupsen/logrus"
	"io"
	"music/internal/readers"
)

type Music struct {
	scenes.Base
	factory   *ObjectFactory
	processor *AudioProcessor
}

func NewMusic(base scenes.Base, f *ObjectFactory) *Music {
	return &Music{
		Base:    base,
		factory: f,
	}
}

func (l *Music) Init() {
	_ = l.factory.NewLine(0, 0, 1, 1)

	reader, err := readers.NewWavReader("./resources/test.wav")
	if err != nil {
		logrus.WithError(err).Fatal("error on creating wav reader")
	}

	l.processor, err = NewAudioProcessor(reader)
	if err != nil {
		logrus.WithError(err).Fatal("error on creating audio processor")
	}
}

// normalization
func mapper(val, vmin, vmax, min, max int) int {
	return (((val - vmin) * (max - min)) / (vmax - vmin)) + min
}

func (l *Music) Update() {
	min, max, part, err := l.processor.Process()
	if err != nil {
		if err == io.EOF {
			return
		}
		logrus.WithError(err).Warn("error on processing")
	}

	x := float32(part)
	y1 := float32(mapper(min, -31250, 31250, 0, 400))
	y2 := float32(mapper(max, -31250, 31250, 0, 400))
	line := l.factory.NewLine(x, y1, x, y2)

	l.DrawObjects.Put(line)
}
