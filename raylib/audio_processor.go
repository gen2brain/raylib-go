//go:build !cgo || nocgo || !noaudio
// +build !cgo nocgo !noaudio

package rl

import (
	"reflect"
	"sync"
	"unsafe"
)

const (
	// audioProcessorSlots is the number of stream processors that can be attached at the same time
	audioProcessorSlots = 64
	// audioDeviceChannels is AUDIO_DEVICE_CHANNELS, processors receive frames in the mixing format
	audioDeviceChannels = 2
)

type audioStreamProcessor struct {
	buffer    *AudioBuffer
	value     uintptr
	code      uintptr
	seq       uint64
	processor AudioCallback
}

var (
	audioProcessorsMutex  sync.RWMutex
	audioProcessorSeq     uint64
	audioStreamProcessors [audioProcessorSlots]audioStreamProcessor
	audioMixedProcessors  []AudioCallback
)

// audioCallbackValue identifies one func value (closure or method value), audioCallbackCode its function
func audioCallbackValue(processor AudioCallback) uintptr {
	return uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&processor)))
}

func audioCallbackCode(processor AudioCallback) uintptr {
	return reflect.ValueOf(processor).Pointer()
}

// acquireAudioStreamProcessor reserves a trampoline slot for processor on stream, -1 when all are in use
func acquireAudioStreamProcessor(stream AudioStream, processor AudioCallback) int {
	audioProcessorsMutex.Lock()
	defer audioProcessorsMutex.Unlock()

	for i := range audioStreamProcessors {
		if audioStreamProcessors[i].processor == nil {
			audioProcessorSeq++
			audioStreamProcessors[i] = audioStreamProcessor{stream.Buffer, audioCallbackValue(processor), audioCallbackCode(processor), audioProcessorSeq, processor}
			return i
		}
	}

	TraceLog(LogWarning, "AUDIO: All %d stream processor slots are in use", audioProcessorSlots)

	return -1
}

// findAudioStreamProcessor returns the slot holding processor on stream, matching the same func value first
// and otherwise the most recent attachment of the same function, -1 when it is not attached
func findAudioStreamProcessor(stream AudioStream, processor AudioCallback) int {
	audioProcessorsMutex.RLock()
	defer audioProcessorsMutex.RUnlock()

	value, code := audioCallbackValue(processor), audioCallbackCode(processor)
	byValue, byCode := -1, -1
	for i := range audioStreamProcessors {
		s := &audioStreamProcessors[i]
		if s.processor == nil || s.buffer != stream.Buffer {
			continue
		}
		if s.value == value && (byValue < 0 || s.seq > audioStreamProcessors[byValue].seq) {
			byValue = i
		}
		if s.code == code && (byCode < 0 || s.seq > audioStreamProcessors[byCode].seq) {
			byCode = i
		}
	}
	if byValue >= 0 {
		return byValue
	}

	return byCode
}

func releaseAudioStreamProcessor(slot int) {
	audioProcessorsMutex.Lock()
	audioStreamProcessors[slot] = audioStreamProcessor{}
	audioProcessorsMutex.Unlock()
}

func callAudioStreamProcessor(slot int, data unsafe.Pointer, frames int) {
	audioProcessorsMutex.RLock()
	processor := audioStreamProcessors[slot].processor
	audioProcessorsMutex.RUnlock()

	if processor != nil {
		processor(unsafe.Slice((*float32)(data), frames*audioDeviceChannels), frames)
	}
}

// addAudioMixedProcessor appends processor and reports whether it is the first one
func addAudioMixedProcessor(processor AudioCallback) bool {
	audioProcessorsMutex.Lock()
	defer audioProcessorsMutex.Unlock()

	audioMixedProcessors = append(audioMixedProcessors, processor)

	return len(audioMixedProcessors) == 1
}

// removeAudioMixedProcessor removes processor, matched like findAudioStreamProcessor, and reports whether none are left
func removeAudioMixedProcessor(processor AudioCallback) bool {
	audioProcessorsMutex.Lock()
	defer audioProcessorsMutex.Unlock()

	value, code := audioCallbackValue(processor), audioCallbackCode(processor)
	found := -1
	for i := len(audioMixedProcessors) - 1; i >= 0; i-- {
		if audioCallbackValue(audioMixedProcessors[i]) == value {
			found = i
			break
		}
		if audioCallbackCode(audioMixedProcessors[i]) == code && found < 0 {
			found = i
		}
	}
	if found >= 0 {
		audioMixedProcessors = append(audioMixedProcessors[:found:found], audioMixedProcessors[found+1:]...)
	}

	return len(audioMixedProcessors) == 0
}

func callAudioMixedProcessors(data unsafe.Pointer, frames int) {
	audioProcessorsMutex.RLock()
	processors := audioMixedProcessors
	audioProcessorsMutex.RUnlock()

	samples := unsafe.Slice((*float32)(data), frames*audioDeviceChannels)
	for _, processor := range processors {
		processor(samples, frames)
	}
}
