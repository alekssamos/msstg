package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const Mp3First4bytes = "\xff\xf3d\xc4"

// this is an mp3 file or not?
func IsMp3(filename string) bool {
	var err error
	f, err := os.Open(filename)
	LogError(err)
	if err != nil {
		return false
	}
	defer f.Close()
	buff := make([]byte, 4)
	_, err = io.ReadAtLeast(f, buff, 4)
	LogError(err)
	if err != nil {
		return false
	}
	if !bytes.Equal(buff, []byte(Mp3First4bytes)) {
		return false
	}
	return true
}

// Get Duration in seconds
func DurationMp3(filename string, bitrate int) (int, error) {
	var durationSeconds int
	if !IsMp3(filename) {
		return 0, fmt.Errorf("this file is not an mp3")
	}
	stat, err := os.Stat(filename)
	LogError(err)
	if err != nil {
		return 0, err
	}
	durationSeconds = int(stat.Size() * int64(8) / int64(bitrate))
	return durationSeconds, nil
}
