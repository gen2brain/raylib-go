// +build !noaudio

package rl

/*
#include "external/stb_vorbis.c"

#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

// cptr returns C pointer
func (w *Wave) cptr() *C.Wave {
	return (*C.Wave)(unsafe.Pointer(w))
}

func (s *Sound) cptr() *C.Sound {
	return (*C.Sound)(unsafe.Pointer(s))
}

// cptr returns C pointer
func (a *AudioStream) cptr() *C.AudioStream {
	return (*C.AudioStream)(unsafe.Pointer(a))
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
	v := bool(int(ret) == 1)
	return v
}

// SetMasterVolume - Set master volume (listener)
func SetMasterVolume(volume float32) {
	cvolume := (C.float)(volume)
	C.SetMasterVolume(cvolume)
}

// LoadWave - Load wave data from file into RAM
func LoadWave(fileName string) Wave {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadWave(cfileName)
	v := newWaveFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadWaveEx - Load wave data from float array data (32bit)
func LoadWaveEx(data []byte, sampleCount int32, sampleRate int32, sampleSize int32, channels int32) Wave {
	cdata := unsafe.Pointer(&data[0])
	csampleCount := (C.int)(sampleCount)
	csampleRate := (C.int)(sampleRate)
	csampleSize := (C.int)(sampleSize)
	cchannels := (C.int)(channels)
	ret := C.LoadWaveEx(cdata, csampleCount, csampleRate, csampleSize, cchannels)
	v := newWaveFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadSound - Load sound to memory
func LoadSound(fileName string) Sound {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadSound(cfileName)
	v := newSoundFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadSoundFromWave - Load sound to memory from wave data
func LoadSoundFromWave(wave Wave) Sound {
	cwave := wave.cptr()
	ret := C.LoadSoundFromWave(*cwave)
	v := newSoundFromPointer(unsafe.Pointer(&ret))
	return v
}

// UpdateSound - Update sound buffer with new data
func UpdateSound(sound Sound, data []byte, samplesCount int32) {
	csound := sound.cptr()
	cdata := unsafe.Pointer(&data[0])
	csamplesCount := (C.int)(samplesCount)
	C.UpdateSound(*csound, cdata, csamplesCount)
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

// PlaySound - Play a sound
func PlaySound(sound Sound) {
	csound := sound.cptr()
	C.PlaySound(*csound)
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

// StopSound - Stop playing a sound
func StopSound(sound Sound) {
	csound := sound.cptr()
	C.StopSound(*csound)
}

// IsSoundPlaying - Check if a sound is currently playing
func IsSoundPlaying(sound Sound) bool {
	csound := sound.cptr()
	ret := C.IsSoundPlaying(*csound)
	v := bool(int(ret) == 1)
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

// WaveFormat - Convert wave data to desired format
func WaveFormat(wave Wave, sampleRate int32, sampleSize int32, channels int32) {
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

// WaveCrop - Crop a wave to defined samples range
func WaveCrop(wave Wave, initSample int32, finalSample int32) {
	cwave := wave.cptr()
	cinitSample := (C.int)(initSample)
	cfinalSample := (C.int)(finalSample)
	C.WaveCrop(cwave, cinitSample, cfinalSample)
}

// GetWaveData - Get samples data from wave as a floats array
func GetWaveData(wave Wave) []float32 {
	var data []float32
	cwave := wave.cptr()
	ret := C.GetWaveData(*cwave)

	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
	sliceHeader.Cap = int(wave.SampleCount)
	sliceHeader.Len = int(wave.SampleCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(ret))

	return data
}

// LoadMusicStream - Load music stream from file
func LoadMusicStream(fileName string) Music {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadMusicStream(cfileName)
	v := *(*Music)(unsafe.Pointer(&ret))
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

// IsMusicPlaying - Check if music is playing
func IsMusicPlaying(music Music) bool {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	ret := C.IsMusicPlaying(cmusic)
	v := bool(int(ret) == 1)
	return v
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

// SetMusicLoopCount - Set music loop count (loop repeats)
// NOTE: If set to -1, means infinite loop
func SetMusicLoopCount(music Music, count int32) {
	cmusic := *(*C.Music)(unsafe.Pointer(&music))
	ccount := (C.int)(count)
	C.SetMusicLoopCount(cmusic, ccount)
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

// InitAudioStream - Init audio stream (to stream raw audio pcm data)
func InitAudioStream(sampleRate uint32, sampleSize uint32, channels uint32) AudioStream {
	csampleRate := (C.uint)(sampleRate)
	csampleSize := (C.uint)(sampleSize)
	cchannels := (C.uint)(channels)
	ret := C.InitAudioStream(csampleRate, csampleSize, cchannels)
	v := newAudioStreamFromPointer(unsafe.Pointer(&ret))
	return v
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream(stream AudioStream, data []float32, samplesCount int32) {
	cstream := stream.cptr()
	cdata := unsafe.Pointer(&data[0])
	csamplesCount := (C.int)(samplesCount)
	C.UpdateAudioStream(*cstream, cdata, csamplesCount)
}

// CloseAudioStream - Close audio stream and free memory
func CloseAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.CloseAudioStream(*cstream)
}

// IsAudioBufferProcessed - Check if any audio stream buffers requires refill
func IsAudioBufferProcessed(stream AudioStream) bool {
	cstream := stream.cptr()
	ret := C.IsAudioBufferProcessed(*cstream)
	v := bool(int(ret) == 1)
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

// StopAudioStream - Stop audio stream
func StopAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.StopAudioStream(*cstream)
}
