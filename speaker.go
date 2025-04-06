package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/alekssamos/edge-tts-go/edgeTTS"
)

const BITRATE = 48000

func ConvEdgeTtsVal(v int, suffix string) (string, error) {
	suffixes := []string{"%", "Hz"}
	foundsuffix := false
	for _, s := range suffixes {
		if s == suffix {
			foundsuffix = true
		}
	}
	if !foundsuffix {
		return "", fmt.Errorf("incorrect suffix. Expected: %s, you passed: %s", strings.Join(suffixes, ", "), suffix)
	}
	sign := "+"
	if v < 0 {
		sign = ""
	}
	if v > 100 || v < -100 {
		return "", fmt.Errorf("the value should be between -100 and +100")
	}
	return fmt.Sprintf("%s%d%s", sign, v, suffix), nil
}

func Speak(text, voice, writeMedia string, rate, pitch, volume int) bool {
	var err error
	const MAXRETRY = 3
	if text == "" {
		log.Println("empty text")
		return false
	}
	crate, err := ConvEdgeTtsVal(rate, "%")
	if err != nil {
		log.Println("rate: " + err.Error())
	}
	cpitch, err := ConvEdgeTtsVal(pitch, "Hz")
	if err != nil {
		log.Println("pitch: " + err.Error())
	}
	cvolume, err := ConvEdgeTtsVal(volume, "%")
	if err != nil {
		log.Println("volume: " + err.Error())
	}
	for i := 1; i <= MAXRETRY; i++ {
		if i > 1 {
			mu.Lock()
			defer mu.Unlock()
		}
		err := edgeTTS.NewTTS(writeMedia).AddText(text, voice, crate, cvolume, cpitch).Speak()
		if err != nil {
			log.Printf("Speak error: %s\n", err.Error())
			continue
		}
		result := IsMp3(writeMedia)
		if result {
			return true
		}
		if i == MAXRETRY {
			log.Println("couldn't voice the text")
		}
	}
	return false
}
