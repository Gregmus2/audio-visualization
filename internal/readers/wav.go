package readers

import (
	"github.com/pkg/errors"
	"github.com/youpy/go-wav"
	"math"
	"os"
)

type Wav struct {
	samples     []int
	dec         *wav.Reader
	sampleRate  int
	numChannels int
}

func NewWavReader(filename string) (Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dec := wav.NewReader(file)

	format, err := dec.Format()
	if err != nil {
		return nil, err
	}
	if format.AudioFormat != wav.AudioFormatPCM {
		return nil, errors.New("audio format is invalid")
	}

	samples, err := wav.NewReader(file).ReadSamples(math.MaxInt32)
	if err != nil {
		return nil, errors.Wrap(err, "error on reading wav samples")
	}

	normalizedSamples := make([]int, len(samples))
	for i, sample := range samples {
		normalizedSamples[i] = dec.IntValue(sample, 0)
	}

	return &Wav{
		samples:     normalizedSamples,
		dec:         dec,
		sampleRate:  int(format.SampleRate),
		numChannels: int(format.NumChannels),
	}, nil
}

func (r *Wav) Samples() []int {
	return r.samples
}

func (r *Wav) SampleRate() int {
	return r.sampleRate
}

func (r *Wav) NumChannels() int {
	return r.numChannels
}

func (r *Wav) Read(size int) ([]byte, error) {
	data := make([]byte, size)
	_, err := r.dec.Read(data)

	return data, err
}
