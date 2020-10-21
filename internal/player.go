package internal

import (
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
	"io"
	"music/internal/readers"
)

const bufferSize = 1024

type AudioProcessor struct {
	player *oto.Player
	offset int
	reader readers.Reader
}

func NewAudioProcessor(reader readers.Reader) (*AudioProcessor, error) {
	ctx, err := oto.NewContext(reader.SampleRate(), reader.NumChannels(), 2, bufferSize)
	if err != nil || ctx == nil {
		return nil, errors.Wrap(err, "error on creating oto context")
	}

	player := ctx.NewPlayer()

	return &AudioProcessor{player: player, reader: reader}, nil
}

func (h *AudioProcessor) Process() (min, max, part int, err error) {
	data, err := h.reader.Read(bufferSize)
	if err == io.EOF {
		if err = h.player.Close(); err != nil {
			return
		}
		return
	}
	if err != nil {
		return
	}

	_, err = h.player.Write(data)
	if err != nil {
		return
	}

	batch := bufferSize / 4
	samples := h.reader.Samples()
	if len(samples) < h.offset+batch {
		return 0, 0, 0, errors.New("wrong samle length")
	}

	samplePart := samples[h.offset : h.offset+batch]
	min = samples[0]
	max = min
	for _, sml := range samplePart {
		if sml < min {
			min = sml
		}
		if sml > max {
			max = sml
		}
	}

	h.offset += batch
	part = h.offset / batch

	return
}
