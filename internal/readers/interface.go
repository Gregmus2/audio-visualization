package readers

type Reader interface {
	Samples() []int
	SampleRate() int
	NumChannels() int
	Read(size int) ([]byte, error)
}
