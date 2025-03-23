package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed voices_list.json
var vj []byte

type Voice struct {
	Name           string `json:"Name"`
	ShortName      string `json:"ShortName"`
	Gender         string `json:"Gender"`
	Locale         string `json:"Locale"`
	SuggestedCodec string `json:"SuggestedCodec"`
	FriendlyName   string `json:"FriendlyName"`
	Status         string `json:"Status"`
	Language       string
	VoiceTag       VoiceTag `json:"VoiceTag"`
}

type VoiceTag struct {
	ContentCategories  []string `json:"ContentCategories"`
	VoicePersonalities []string `json:"VoicePersonalities"`
}

// returns only voice name without locale and "Neural" word. E.G: Adri, Willem
func (v Voice) OnlyName() string {
	ss := strings.Split(v.ShortName, "-")
	return strings.TrimSuffix(ss[2], "Neural")
}

func (v Voice) Country() string {
	return strings.Split(strings.Split(v.FriendlyName, " - ")[1], " ")[0]
}

func ListVoices() ([]Voice, error) {
	// Parse the JSON response.
	var voices []Voice
	err := json.Unmarshal(vj, &voices)
	LogError(err)
	if err != nil {
		return nil, err
	}

	return voices, nil
}

func FindVoices(query string) ([]Voice, error) {
	voices, err := ListVoices()
	var foundVoices []Voice
	LogError(err)
	if err != nil {
		return foundVoices, err
	}
	for _, v := range voices {
		if strings.Contains(strings.ToLower(fmt.Sprintf("%s %s", v.ShortName, v.FriendlyName)), strings.ToLower(query)) {
			foundVoices = append(foundVoices, v)
		}
	}
	if len(foundVoices) > 0 {
		return foundVoices, nil
	} else {
		return foundVoices, fmt.Errorf("the query %q yielded no results", query)
	}
}

func AllLocales() []string {
	sLocales := " "
	voices, _ := ListVoices()
	for _, v := range voices {
		if strings.Contains(sLocales, v.Locale) {
			continue
		}
		sLocales = sLocales + " " + v.Locale
	}
	return strings.Fields(sLocales)
}

func AllCountries() []string {
	sCountries := " "
	voices, _ := ListVoices()
	for _, v := range voices {
		country := strings.Split(strings.Split(v.FriendlyName, " - ")[1], " ")[0]
		if strings.Contains(sCountries, country) {
			continue
		}
		sCountries = sCountries + " " + country
	}
	return strings.Fields(sCountries)
}
