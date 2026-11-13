//go:build !noaudio && !nocgo
// +build !noaudio,!nocgo

package rl

/*
#include "raylib.h"
#include <stdlib.h>

#if defined(_WIN32) || defined(WIN32) || defined(__WIN32__) || defined(__NT__)
    __declspec(dllexport)
#endif
extern void internalAudioStreamCallbackGo(void *, int);

static void audioStreamWrapperCallback(void *data, unsigned int frames) {
	internalAudioStreamCallbackGo(data, frames);
}

static void setAudioStreamCallbackWrapper(AudioStream stream) {
	SetAudioStreamCallback(stream, audioStreamWrapperCallback);
}

#if defined(_WIN32) || defined(WIN32) || defined(__WIN32__) || defined(__NT__)
    __declspec(dllexport)
#endif
extern void internalAudioMixedProcessorGo(void *, int);

static void audioMixedProcessorCallback(void *data, unsigned int frames) {
	internalAudioMixedProcessorGo(data, frames);
}

static void setAudioMixedProcessorCallbackWrapper() {
	AttachAudioMixedProcessor(audioMixedProcessorCallback);
}

static void unsetAudioMixedProcessorCallbackWrapper() {
	DetachAudioMixedProcessor(audioMixedProcessorCallback);
}

#if defined(_WIN32) || defined(WIN32) || defined(__WIN32__) || defined(__NT__)
    __declspec(dllexport)
#endif
extern void internalAudioStreamProcessorGo(int, void *, int);

#define AUDIO_STREAM_PROCESSOR(n) static void audioStreamProcessor##n(void *data, unsigned int frames) { internalAudioStreamProcessorGo(n, data, frames); }
#define AUDIO_STREAM_PROCESSOR8(n) AUDIO_STREAM_PROCESSOR(n##0) AUDIO_STREAM_PROCESSOR(n##1) AUDIO_STREAM_PROCESSOR(n##2) AUDIO_STREAM_PROCESSOR(n##3) \
	AUDIO_STREAM_PROCESSOR(n##4) AUDIO_STREAM_PROCESSOR(n##5) AUDIO_STREAM_PROCESSOR(n##6) AUDIO_STREAM_PROCESSOR(n##7)
AUDIO_STREAM_PROCESSOR8(0) AUDIO_STREAM_PROCESSOR8(1) AUDIO_STREAM_PROCESSOR8(2) AUDIO_STREAM_PROCESSOR8(3)
AUDIO_STREAM_PROCESSOR8(4) AUDIO_STREAM_PROCESSOR8(5) AUDIO_STREAM_PROCESSOR8(6) AUDIO_STREAM_PROCESSOR8(7)

#define AUDIO_STREAM_PROCESSOR_REF8(n) audioStreamProcessor##n##0, audioStreamProcessor##n##1, audioStreamProcessor##n##2, audioStreamProcessor##n##3, \
	audioStreamProcessor##n##4, audioStreamProcessor##n##5, audioStreamProcessor##n##6, audioStreamProcessor##n##7
static AudioCallback audioStreamProcessors[64] = {
	AUDIO_STREAM_PROCESSOR_REF8(0), AUDIO_STREAM_PROCESSOR_REF8(1), AUDIO_STREAM_PROCESSOR_REF8(2), AUDIO_STREAM_PROCESSOR_REF8(3),
	AUDIO_STREAM_PROCESSOR_REF8(4), AUDIO_STREAM_PROCESSOR_REF8(5), AUDIO_STREAM_PROCESSOR_REF8(6), AUDIO_STREAM_PROCESSOR_REF8(7),
};

static void attachAudioStreamProcessorSlot(AudioStream stream, int slot) {
	AttachAudioStreamProcessor(stream, audioStreamProcessors[slot]);
}

static void detachAudioStreamProcessorSlot(AudioStream stream, int slot) {
	DetachAudioStreamProcessor(stream, audioStreamProcessors[slot]);
}
*/
import "C"
import (
	"unsafe"
)

var (
	internalAudioStreamCallback AudioCallback
	internalAudioStreamChannels uint32 = 2
)

// SetAudioStreamCallback - Audio thread callback to request new data
func SetAudioStreamCallback(stream AudioStream, callback AudioCallback) {
	internalAudioStreamCallback = callback
	internalAudioStreamChannels = stream.Channels
	C.setAudioStreamCallbackWrapper(*stream.cptr())
}

//export internalAudioStreamCallbackGo
func internalAudioStreamCallbackGo(data unsafe.Pointer, frames C.int) {
	if internalAudioStreamCallback != nil {
		// Buffer holds frames*channels samples, interleaved per stream channels.
		channels := int(internalAudioStreamChannels)
		if channels < 1 {
			channels = 1
		}
		internalAudioStreamCallback(unsafe.Slice((*float32)(data), int(frames)*channels), int(frames))
	}
}

// AttachAudioStreamProcessor - Attach audio stream processor to stream, receives frames x 2 samples as 'float' (stereo)
func AttachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	slot := acquireAudioStreamProcessor(stream, processor)
	if slot < 0 {
		return
	}

	C.attachAudioStreamProcessorSlot(*stream.cptr(), C.int(slot))
}

// DetachAudioStreamProcessor - Detach audio stream processor from stream
func DetachAudioStreamProcessor(stream AudioStream, processor AudioCallback) {
	slot := findAudioStreamProcessor(stream, processor)
	if slot < 0 {
		return
	}

	C.detachAudioStreamProcessorSlot(*stream.cptr(), C.int(slot))
	releaseAudioStreamProcessor(slot)
}

//export internalAudioStreamProcessorGo
func internalAudioStreamProcessorGo(slot C.int, data unsafe.Pointer, frames C.int) {
	callAudioStreamProcessor(int(slot), data, int(frames))
}

// AttachAudioMixedProcessor - Attach audio stream processor to the entire audio pipeline, receives frames x 2 samples as 'float' (stereo)
func AttachAudioMixedProcessor(processor AudioCallback) {
	if addAudioMixedProcessor(processor) {
		C.setAudioMixedProcessorCallbackWrapper()
	}
}

// DetachAudioMixedProcessor - Detach audio stream processor from the entire audio pipeline
func DetachAudioMixedProcessor(processor AudioCallback) {
	if removeAudioMixedProcessor(processor) {
		C.unsetAudioMixedProcessorCallbackWrapper()
	}
}

//export internalAudioMixedProcessorGo
func internalAudioMixedProcessorGo(data unsafe.Pointer, frames C.int) {
	callAudioMixedProcessors(data, int(frames))
}

// newWaveFromPointer - Returns new Wave from pointer
func newWaveFromPointer(ptr unsafe.Pointer) Wave {
	return *(*Wave)(ptr)
}

// cptr returns C pointer
func (w *Wave) cptr() *C.Wave {
	return (*C.Wave)(unsafe.Pointer(w))
}

// newSoundFromPointer - Returns new Sound from pointer
func newSoundFromPointer(ptr unsafe.Pointer) Sound {
	return *(*Sound)(ptr)
}

func (s *Sound) cptr() *C.Sound {
	return (*C.Sound)(unsafe.Pointer(s))
}

// newAudioStreamFromPointer - Returns new AudioStream from pointer
func newAudioStreamFromPointer(ptr unsafe.Pointer) AudioStream {
	return *(*AudioStream)(ptr)
}

// cptr returns C pointer
func (a *AudioStream) cptr() *C.AudioStream {
	return (*C.AudioStream)(unsafe.Pointer(a))
}

// newMusicFromPointer - Returns new Music from pointer
func newMusicFromPointer(ptr unsafe.Pointer) Music {
	return *(*Music)(ptr)
}

// InitAudioDevice - Initialize audio device and context
func InitAudioDevice() {
	C.InitAudioDevice()
}

// CloseAudioDevice - Close the audio device and context
func CloseAudioDevice() {
	C.CloseAudioDevice()
}

// IsAudioDeviceReady - Check if audio device has been initialized successfully
func IsAudioDeviceReady() bool {
	ret := C.IsAudioDeviceReady()
	v := bool(ret)
	return v
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {
	cvolume := (C.float)(volume)
	C.SetMasterVolume(cvolume)
}

// GetMasterVolume - Get master volume (listener)
func GetMasterVolume() float32 {
	ret := C.GetMasterVolume()
	v := float32(ret)
	return v
}

// LoadWave - Load wave data from file
func LoadWave(fileName string) Wave {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadWave(cfileName)
	v := newWaveFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadWaveFromMemory - Load wave from memory buffer, fileType refers to extension: i.e. '.wav'
func LoadWaveFromMemory(fileType string, fileData []byte, dataSize int32) Wave {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	ret := C.LoadWaveFromMemory(cfileType, cfileData, cdataSize)
	v := newWaveFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsWaveValid - Checks if wave data is valid (data loaded and parameters)
func IsWaveValid(wave Wave) bool {
	cwave := wave.cptr()
	ret := C.IsWaveValid(*cwave)
	v := bool(ret)
	return v
}

// LoadSound - Load sound from file
func LoadSound(fileName string) Sound {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadSound(cfileName)
	v := newSoundFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadSoundFromWave - Load sound from wave data
func LoadSoundFromWave(wave Wave) Sound {
	cwave := wave.cptr()
	ret := C.LoadSoundFromWave(*cwave)
	v := newSoundFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadSoundAlias - Create a new sound that shares the same sample data as the source sound, does not own the sound data
func LoadSoundAlias(source Sound) Sound {
	csound := source.cptr()
	ret := C.LoadSoundAlias(*csound)
	v := newSoundFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsSoundValid - Checks if a sound is valid (data loaded and buffers initialized)
func IsSoundValid(sound Sound) bool {
	csound := sound.cptr()
	ret := C.IsSoundValid(*csound)
	v := bool(ret)
	return v
}

// UpdateSound - Update sound buffer with new data (default data format: 32 bit float, stereo)
func UpdateSound(sound Sound, data []byte, sampleCount int32) {
	csound := sound.cptr()
	cdata := unsafe.Pointer(&data[0])
	csampleCount := (C.int)(sampleCount)
	C.UpdateSound(*csound, cdata, csampleCount)
}

// UnloadWave - Unload wave data
func UnloadWave(wave Wave) {
	cwave := wave.cptr()
	C.UnloadWave(*cwave)
}

// UnloadSound - Unload sound
func UnloadSound(sound Sound) {
	csound := sound.cptr()
	C.UnloadSound(*csound)
}

// UnloadSoundAlias - Unload a sound alias (does not deallocate sample data)
func UnloadSoundAlias(alias Sound) {
	csound := alias.cptr()
	C.UnloadSoundAlias(*csound)
}

// ExportWave - Export wave data to file, returns true on success
func ExportWave(wave Wave, fileName string) bool {
	cwave := wave.cptr()
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.ExportWave(*cwave, cfileName)
	v := bool(ret)
	return v
}

// PlaySound - Play a sound
func PlaySound(sound Sound) {
	csound := sound.cptr()
	C.PlaySound(*csound)
}

// StopSound - Stop playing a sound
func StopSound(sound Sound) {
	csound := sound.cptr()
	C.StopSound(*csound)
}

// PauseSound - Pause a sound
func PauseSound(sound Sound) {
	csound := sound.cptr()
	C.PauseSound(*csound)
}

// ResumeSound - Resume a paused sound
func ResumeSound(sound Sound) {
	csound := sound.cptr()
	C.ResumeSound(*csound)
}

// IsSoundPlaying - Check if a sound is currently playing
func IsSoundPlaying(sound Sound) bool {
	csound := sound.cptr()
	ret := C.IsSoundPlaying(*csound)
	v := bool(ret)
	return v
}

// SetSoundVolume - Set volume for a sound (1.0 is max level)
func SetSoundVolume(sound Sound, volume float32) {
	csound := sound.cptr()
	cvolume := (C.float)(volume)
	C.SetSoundVolume(*csound, cvolume)
}

// SetSoundPitch - Set pitch for a sound (1.0 is base level)
func SetSoundPitch(sound Sound, pitch float32) {
	csound := sound.cptr()
	cpitch := (C.float)(pitch)
	C.SetSoundPitch(*csound, cpitch)
}

// SetSoundPan - Set pan for a sound (-1.0 left, 0.0 center, 1.0 right)
func SetSoundPan(sound Sound, pan float32) {
	csound := sound.cptr()
	cpan := (C.float)(pan)
	C.SetSoundPan(*csound, cpan)
}

// WaveFormat - Convert wave data to desired format
func WaveFormat(wave *Wave, sampleRate int32, sampleSize int32, channels int32) {
	cwave := wave.cptr()
	csampleRate := (C.int)(sampleRate)
	csampleSize := (C.int)(sampleSize)
	cchannels := (C.int)(channels)
	C.WaveFormat(cwave, csampleRate, csampleSize, cchannels)
}

// WaveCopy - Copy a wave to a new wave
func WaveCopy(wave Wave) Wave {
	cwave := wave.cptr()
	ret := C.WaveCopy(*cwave)
	v := newWaveFromPointer(unsafe.Pointer(&ret))
	return v
}

// WaveCrop - Crop a wave to defined frames range
func WaveCrop(wave *Wave, initFrame int32, finalFrame int32) {
	cwave := wave.cptr()
	cinitFrame := (C.int)(initFrame)
	cfinalFrame := (C.int)(finalFrame)
	C.WaveCrop(cwave, cinitFrame, cfinalFrame)
}

// LoadWaveSamples - Load samples data from wave as a 32bit float data array
func LoadWaveSamples(wave Wave) []float32 {
	cwave := wave.cptr()
	ret := C.LoadWaveSamples(*cwave)
	v := unsafe.Slice((*float32)(unsafe.Pointer(ret)), wave.FrameCount*wave.Channels)
	return v
}

// UnloadWaveSamples - Unload samples data loaded with LoadWaveSamples()
func UnloadWaveSamples(samples []float32) {
	C.UnloadWaveSamples((*C.float)(unsafe.Pointer(&samples[0])))
}

// LoadMusicStream - Load music stream from file
func LoadMusicStream(fileName string) Music {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadMusicStream(cfileName)
	v := newMusicFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadMusicStreamFromMemory - Load music stream from data
func LoadMusicStreamFromMemory(fileType string, data []byte, dataSize int32) Music {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cdata := (*C.uchar)(unsafe.Pointer(&data[0]))
	cdataSize := (C.int)(dataSize)
	ret := C.LoadMusicStreamFromMemory(cfileType, cdata, cdataSize)
	v := newMusicFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsMusicValid - Checks if a music stream is valid (context and buffers initialized)
func IsMusicValid(music Music) bool {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	ret := C.IsMusicValid(cmusic)
	v := bool(ret)
	return v
}

// UnloadMusicStream - Unload music stream
func UnloadMusicStream(music Music) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	C.UnloadMusicStream(cmusic)
}

// PlayMusicStream - Start music playing
func PlayMusicStream(music Music) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	C.PlayMusicStream(cmusic)
}

// IsMusicStreamPlaying - Check if music is playing
func IsMusicStreamPlaying(music Music) bool {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	ret := C.IsMusicStreamPlaying(cmusic)
	v := bool(ret)
	return v
}

// UpdateMusicStream - Updates buffers for music streaming
func UpdateMusicStream(music Music) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	C.UpdateMusicStream(cmusic)
}

// StopMusicStream - Stop music playing
func StopMusicStream(music Music) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	C.StopMusicStream(cmusic)
}

// PauseMusicStream - Pause music playing
func PauseMusicStream(music Music) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	C.PauseMusicStream(cmusic)
}

// ResumeMusicStream - Resume playing paused music
func ResumeMusicStream(music Music) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	C.ResumeMusicStream(cmusic)
}

// SeekMusicStream - Seek music to a position (in seconds)
func SeekMusicStream(music Music, position float32) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	cposition := (C.float)(position)
	C.SeekMusicStream(cmusic, cposition)
}

// SetMusicVolume - Set volume for music (1.0 is max level)
func SetMusicVolume(music Music, volume float32) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	cvolume := (C.float)(volume)
	C.SetMusicVolume(cmusic, cvolume)
}

// SetMusicPitch - Set pitch for a music (1.0 is base level)
func SetMusicPitch(music Music, pitch float32) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	cpitch := (C.float)(pitch)
	C.SetMusicPitch(cmusic, cpitch)
}

// SetMusicPan - Set pan for a music (-1.0 left, 0.0 center, 1.0 right)
func SetMusicPan(music Music, pan float32) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	cpan := (C.float)(pan)
	C.SetMusicPan(cmusic, cpan)
}

// GetMusicTimeLength - Get music time length (in seconds)
func GetMusicTimeLength(music Music) float32 {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	ret := C.GetMusicTimeLength(cmusic)
	v := (float32)(ret)
	return v
}

// GetMusicTimePlayed - Get current music time played (in seconds)
func GetMusicTimePlayed(music Music) float32 {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	ret := C.GetMusicTimePlayed(cmusic)
	v := (float32)(ret)
	return v
}

// LoadAudioStream - Load audio stream (to stream raw audio pcm data)
func LoadAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {
	csampleRate := (C.uint)(sampleRate)
	csampleSize := (C.uint)(sampleSize)
	cchannels := (C.uint)(channels)
	ret := C.LoadAudioStream(csampleRate, csampleSize, cchannels)
	v := newAudioStreamFromPointer(unsafe.Pointer(&ret))
	return v
}

// IsAudioStreamValid - Checks if an audio stream is valid (buffers initialized)
func IsAudioStreamValid(stream AudioStream) bool {
	cstream := stream.cptr()
	ret := C.IsAudioStreamValid(*cstream)
	v := bool(ret)
	return v
}

// UnloadAudioStream - Unload audio stream and free memory
func UnloadAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.UnloadAudioStream(*cstream)
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream[T float32 | int16](stream AudioStream, data []T) {
	var cdata unsafe.Pointer
	switch d := any(data).(type) {
	case []float32:
		cdata = unsafe.Pointer(&d[0])
	case []int16:
		cdata = unsafe.Pointer(&d[0])
	}
	cframeCount := (C.int)(len(data) / int(max(stream.Channels, 1)))
	cstream := stream.cptr()
	C.UpdateAudioStream(*cstream, cdata, cframeCount)
}

// IsAudioStreamProcessed - Check if any audio stream buffers requires refill
func IsAudioStreamProcessed(stream AudioStream) bool {
	cstream := stream.cptr()
	ret := C.IsAudioStreamProcessed(*cstream)
	v := bool(ret)
	return v
}

// PlayAudioStream - Play audio stream
func PlayAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.PlayAudioStream(*cstream)
}

// PauseAudioStream - Pause audio stream
func PauseAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.PauseAudioStream(*cstream)
}

// ResumeAudioStream - Resume audio stream
func ResumeAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.ResumeAudioStream(*cstream)
}

// IsAudioStreamPlaying - Check if audio stream is playing
func IsAudioStreamPlaying(stream AudioStream) bool {
	cstream := stream.cptr()
	ret := C.IsAudioStreamPlaying(*cstream)
	v := bool(ret)
	return v
}

// StopAudioStream - Stop audio stream
func StopAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.StopAudioStream(*cstream)
}

// SetAudioStreamVolume - Set volume for audio stream (1.0 is max level)
func SetAudioStreamVolume(stream AudioStream, volume float32) {
	cstream := stream.cptr()
	cvolume := (C.float)(volume)
	C.SetAudioStreamVolume(*cstream, cvolume)
}

// SetAudioStreamPitch - Set pitch for audio stream (1.0 is base level)
func SetAudioStreamPitch(stream AudioStream, pitch float32) {
	cstream := stream.cptr()
	cpitch := (C.float)(pitch)
	C.SetAudioStreamPitch(*cstream, cpitch)
}

// SetAudioStreamPan - Set pan for audio stream (-1.0 to 1.0 range, 0.0 is centered)
func SetAudioStreamPan(stream AudioStream, pan float32) {
	cstream := stream.cptr()
	cpan := (C.float)(pan)
	C.SetAudioStreamPan(*cstream, cpan)
}

// SetAudioStreamBufferSizeDefault - Default size for new audio streams
func SetAudioStreamBufferSizeDefault(size int32) {
	csize := (C.int)(size)
	C.SetAudioStreamBufferSizeDefault(csize)
}
