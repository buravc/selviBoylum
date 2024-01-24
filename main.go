package main

import (
	"log"
	"math"
	"reflect"
	"selviBoylum/song"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"

var selviBoylum = generate_selviboylum()

var prevIndex = 0

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	if prevIndex >= len(selviBoylum) {
		return
	}
	n := int(length)

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(stream)),
		Len:  n,
		Cap:  n,
	}

	buf := *(*[]C.float)(unsafe.Pointer(&hdr))

	highBound := int(math.Min(float64(prevIndex+(n/4)), float64(len(selviBoylum)-1)))
	samples := selviBoylum[prevIndex:highBound]

	for i := 0; i < len(samples); i++ {
		smp := samples[i]
		buf[i] = C.float(smp)
	}
	prevIndex += n / 4
}

func generate_selviboylum() [song.SAMPLES_PER_SECOND * 32]float32 {
	songBuilder := song.New()

	songBuilder.
		AddNote(song.NOTE_E, 0, 1.0/4).
		AddSilence(1.0/4).
		AddNote(song.NOTE_D, 1, 1.0/8).
		AddNote(song.NOTE_C_S, 1, 1.0/8).
		AddNote(song.NOTE_D, 1, 1.0/8).
		AddNote(song.NOTE_B, 0, 1.0/4).
		AddSilence(1.0/2).
		// part 2
		AddNote(song.NOTE_B, 0, 1.0/8).
		AddNote(song.NOTE_G, 0, 1.0/8).
		AddNote(song.NOTE_B, 0, 1.0/16).
		AddNote(song.NOTE_A, 0, 1.0/16).
		AddNote(song.NOTE_A, 0, 1.0/4).
		AddSilence(1.0/4).
		// part 3
		AddNote(song.NOTE_A, 0, 1.0/16).
		AddNote(song.NOTE_B, 0, 1.0/16).
		AddNote(song.NOTE_A, 0, 1.0/8).
		AddNote(song.NOTE_F*song.SEMITONE, 0, 1.0/8).
		AddNote(song.NOTE_G, 0, 1.0/4).
		// part 4
		AddNote(song.NOTE_F*song.SEMITONE, 0, 1.0/8).
		AddNote(song.NOTE_G, 0, 1.0/16).
		AddNote(song.NOTE_A, 0, 1.0/16).
		AddNote(song.NOTE_F*song.SEMITONE, 0, 1.0/8).
		AddNote(song.NOTE_E, 0, 1.0/8).
		AddNote(song.NOTE_E, 0, 1.0/4)

	return songBuilder.Samples()
}

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	spec := &sdl.AudioSpec{
		Freq:     song.SAMPLES_PER_SECOND,
		Format:   sdl.AUDIO_F32SYS,
		Channels: 1,
		Samples:  song.SAMPLES_PER_SECOND / 4,
		Callback: sdl.AudioCallback(C.SineWave),
	}
	if err := sdl.OpenAudio(spec, nil); err != nil {
		log.Println(err)
		return
	}
	sdl.PauseAudio(false)
	sdl.Delay(15000) // play audio for long enough to understand whether it works
	sdl.CloseAudio()
}
