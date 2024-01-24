package song

import (
	"math"
	"math/rand"
)

const (
	NOTE_A   Note = 440.0
	NOTE_B   Note = 493.9
	NOTE_C   Note = 261.6
	NOTE_C_S Note = 277.2
	NOTE_D   Note = 293.7
	NOTE_E   Note = 329.6
	NOTE_F   Note = 349.2
	NOTE_G   Note = 392.0
)

const SEMITONE = 1.05946

type Note float32

const SAMPLES_PER_SECOND = 48000 // samples per second

type songBuilder struct {
	buffer [SAMPLES_PER_SECOND * 32]float32
	index  int

	noise [SAMPLES_PER_SECOND]float32
}

func generate_guitar(noise []float32, frequency Note, octaveOffset, duration int) []float32 {
	buffer := make([]float32, duration)
	requiredNoiseLen := int(SAMPLES_PER_SECOND / (float64(frequency) * math.Pow(2, float64(octaveOffset))))
	for i := 0; i < requiredNoiseLen; i++ {
		buffer[i] = noise[i]
	}

	for i := requiredNoiseLen; i < duration; i++ {
		buffer[i] = (buffer[i-requiredNoiseLen] + buffer[i-requiredNoiseLen+1]) / 2
	}

	return buffer
}

func New() songBuilder {
	s := songBuilder{}
	for i := range s.noise {
		s.noise[i] = rand.Float32()*2 - 1
	}

	return s
}

func (s *songBuilder) AddNote(note Note, octaveOffset int, length float32) *songBuilder {
	wave := generate_guitar(s.noise[:], note, octaveOffset, int(length*3*SAMPLES_PER_SECOND))
	for i := 0; i < len(wave); i++ {
		s.buffer[s.index] = wave[i]
		s.index++
	}

	return s
}

func (s *songBuilder) AddSilence(length float32) *songBuilder {
	var lastSignal float32 = 0.0
	if s.index > 0 {
		lastSignal = s.buffer[s.index-1]
	}
	for i := 0; i < int(length*SAMPLES_PER_SECOND); i++ {
		s.buffer[s.index] = lastSignal
		s.index++
	}

	return s
}

func (s *songBuilder) Samples() [SAMPLES_PER_SECOND * 32]float32 {
	return s.buffer
}
