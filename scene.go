package main

import (
	"github.com/Gregmus2/simple-engine/scenes"
	"github.com/hajimehoshi/oto"
	"github.com/youpy/go-wav"
	"io"
	"log"
	"math"
	"os"
	"time"
)

const bufferSize = 1024

type Music struct {
	scenes.Base
	factory *ObjectFactory
	dec *wav.Reader
	samples []wav.Sample
	player *oto.Player
	offset int
}

func NewMusic(base scenes.Base, f *ObjectFactory) *Music {
	return &Music{
		Base:    base,
		factory: f,
	}
}

func (l *Music) Init() {
	time.Sleep(4 * time.Second)
	_ = l.factory.NewLine(0, 0, 1, 1)

	var err error

	var file *os.File
	if file, err = os.Open("./test.wav"); err != nil {
		log.Fatal(err)
	}

	l.dec = wav.NewReader(file)

	wavformat, err_rd := l.dec.Format()
	if err_rd != nil {
		panic(err_rd)
	}
	if wavformat.AudioFormat != wav.AudioFormatPCM {
		panic("Audio format is invalid ")
	}
	log.Println(wavformat.SampleRate)
	log.Println(wavformat.NumChannels)

	l.samples, err = wav.NewReader(file).ReadSamples(math.MaxInt32)
	if err != nil {
		panic(err)
	}

	var context *oto.Context
	if context, err = oto.NewContext(int(wavformat.SampleRate), int(wavformat.NumChannels), 2, bufferSize); err != nil {
		log.Fatal(err)
	}

	l.player = context.NewPlayer()
	//w, h := l.Window.GetSize()
}

func mapper(val, vmin, vmax, min, max int) int {
	return (((val-vmin)*(max-min))/(vmax-vmin))+min
}

func (l *Music) Update() {
	var data = make([]byte, bufferSize)

	_, err := l.dec.Read(data)
	if err == io.EOF {
		if err = l.player.Close(); err != nil {
			log.Fatal(err)
		}
		panic(1)
	}
	if err != nil {
		panic(err)
	}
	l.player.Write(data)

	batch := bufferSize / 4
	if len(l.samples) < l.offset+batch {
		panic(1)
	}

	smpls := l.samples[l.offset : l.offset+batch]
	min1 := l.dec.IntValue(smpls[0], 0)
	min2 := l.dec.IntValue(smpls[0], 1)
	max1, max2 := min1, min2
	for _, sml := range smpls {
		val1 := l.dec.IntValue(sml, 0)
		if val1 < min1 {
			min1 = val1
		}
		if val1 > max1 {
			max1 = val1
		}
		val2 := l.dec.IntValue(sml, 0)
		if val2 < min2 {
			min2 = val2
		}
		if val2 > max2 {
			max2 = val2
		}
	}

	x := float32(l.offset / batch)
	line1 := l.factory.NewLine(x, float32(mapper(min1, -31250, 31250, 0, 400)), x, float32(mapper(max1, -31250, 31250, 0, 400)))
	line2 := l.factory.NewLine(x, float32(mapper(min2, -31250, 31250, 400, 800)), x, float32(mapper(max2, -31250, 31250, 400, 800)))

	l.DrawObjects.Put(line1)
	l.DrawObjects.Put(line2)

	l.offset += batch
}
