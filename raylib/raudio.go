//go:build !noaudio
// +build !noaudio

package rl

/*
//#include "external/stb_vorbis.c"

#include "raylib.h"
#include <stdlib.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

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
	v := bool(ret)
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

// LoadWaveFromMemory - Load wave from memory buffer, fileType refers to extension: i.e. ".wav"
func LoadWaveFromMemory(fileType string, fileData []byte, dataSize int32) Wave {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	ret := C.LoadWaveFromMemory(cfileType, cfileData, cdataSize)
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

// ExportWave - Export wave data to file
func ExportWave(wave Wave, fileName string) {
	cwave := wave.cptr()
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	C.ExportWave(*cwave, cfileName)
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

// PlaySoundMulti - Play a sound (using multichannel buffer pool)
func PlaySoundMulti(sound Sound) {
	csound := sound.cptr()
	C.PlaySoundMulti(*csound)
}

// StopSoundMulti - Stop any sound playing (using multichannel buffer pool)
func StopSoundMulti() {
	C.StopSoundMulti()
}

// GetSoundsPlaying - Get number of sounds playing in the multichannel
func GetSoundsPlaying() int {
	ret := C.GetSoundsPlaying()
	v := int(ret)
	return v
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

// SetSoundPan - Set pan for a sound (0.5 is center)
func SetSoundPan(sound Sound, pan float32) {
	csound := sound.cptr()
	cpan := (C.float)(pan)
	C.SetSoundPan(*csound, cpan)
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

// LoadWaveSamples - Get samples data from wave as a floats array
func LoadWaveSamples(wave Wave) []float32 {
	var data []float32
	cwave := wave.cptr()
	ret := C.LoadWaveSamples(*cwave)

	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
	sliceHeader.Cap = int(wave.FrameCount)
	sliceHeader.Len = int(wave.FrameCount)
	sliceHeader.Data = uintptr(unsafe.Pointer(ret))

	return data
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
func LoadMusicStreamFromMemory(fileType string, fileData []byte, dataSize int32) Music {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	ret := C.LoadMusicStreamFromMemory(cfileType, cfileData, cdataSize)
	v := newMusicFromPointer(unsafe.Pointer(&ret))
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

// SetMusicPan - Set pan for a music (0.5 is center)
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

// UnloadAudioStream - Unload audio stream and free memory
func UnloadAudioStream(stream AudioStream) {
	cstream := stream.cptr()
	C.UnloadAudioStream(*cstream)
}

// UpdateAudioStream - Update audio stream buffers with data
func UpdateAudioStream(stream AudioStream, data []float32, samplesCount int32) {
	cstream := stream.cptr()
	cdata := unsafe.Pointer(&data[0])
	csamplesCount := (C.int)(samplesCount)
	C.UpdateAudioStream(*cstream, cdata, csamplesCount)
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

// SetAudioStreamPan - Set pan for audio stream (0.5 is centered)
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
