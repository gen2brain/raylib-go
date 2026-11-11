//go:build !cgo || nocgo || !noaudio

package rl

import (
	"testing"
	"unsafe"
)

type testProcessor struct{ calls, samples int }

func (p *testProcessor) process(data []float32, frames int) {
	p.calls++
	p.samples = len(data)
}

func TestAudioStreamProcessorSlots(t *testing.T) {
	defer func() { audioStreamProcessors = [audioProcessorSlots]audioStreamProcessor{} }()

	stream := AudioStream{Buffer: &AudioBuffer{}}
	other := AudioStream{Buffer: &AudioBuffer{}}

	var a, b testProcessor
	pa, pb := a.process, b.process
	sa := acquireAudioStreamProcessor(stream, pa)
	sb := acquireAudioStreamProcessor(stream, pb)

	if got := findAudioStreamProcessor(stream, pb); got != sb {
		t.Fatalf("find by value = %d, want %d", got, sb)
	}
	if got := findAudioStreamProcessor(other, pa); got != -1 {
		t.Fatalf("find on another stream = %d, want -1", got)
	}

	buf := make([]float32, 8*audioDeviceChannels)
	callAudioStreamProcessor(sa, unsafe.Pointer(&buf[0]), 8)
	if a.calls != 1 || a.samples != 8*audioDeviceChannels || b.calls != 0 {
		t.Fatalf("dispatch: a=%+v b=%+v", a, b)
	}

	releaseAudioStreamProcessor(sa)
	if got := findAudioStreamProcessor(stream, a.process); got != sb {
		t.Fatalf("fresh method value falls back to the attached one = %d, want %d", got, sb)
	}

	releaseAudioStreamProcessor(sb)
	for i := 0; i < audioProcessorSlots; i++ {
		if acquireAudioStreamProcessor(stream, pa) < 0 {
			t.Fatalf("slot %d not available", i)
		}
	}
	if got := acquireAudioStreamProcessor(stream, pa); got != -1 {
		t.Fatalf("acquire beyond %d slots = %d, want -1", audioProcessorSlots, got)
	}
}

func TestAudioMixedProcessors(t *testing.T) {
	defer func() { audioMixedProcessors = nil }()

	var a, b testProcessor
	pa, pb := a.process, b.process
	if !addAudioMixedProcessor(pa) || addAudioMixedProcessor(pb) {
		t.Fatal("first add must report true, second false")
	}

	if removeAudioMixedProcessor(pa) {
		t.Fatal("remove with one left reported empty")
	}

	buf := make([]float32, 4*audioDeviceChannels)
	callAudioMixedProcessors(unsafe.Pointer(&buf[0]), 4)
	if a.calls != 0 || b.calls != 1 || b.samples != 4*audioDeviceChannels {
		t.Fatalf("dispatch after removing a: a=%+v b=%+v", a, b)
	}

	if !removeAudioMixedProcessor(pb) {
		t.Fatal("removing the last must report empty")
	}
}
