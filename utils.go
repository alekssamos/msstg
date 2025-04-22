package main

import "strings"

// russian to english transliteration map
var translitMap = map[string]string{
	"а": "a",
	"б": "b",
	"в": "v",
	"г": "g",
	"д": "d",
	"е": "e",
	"ё": "yo",
	"ж": "zh",
	"з": "z",
	"и": "i",
	"й": "y",
	"к": "k",
	"л": "l",
	"м": "m",
	"н": "n",
	"о": "o",
	"п": "p",
	"р": "r",
	"с": "s",
	"т": "t",
	"у": "u",
	"ф": "f",
	"х": "kh",
	"ц": "ts",
	"ч": "ch",
	"ш": "sh",
	"щ": "shch",
	"ъ": "",
	"ы": "y",
	"ь": "",
	"э": "e",
	"ю": "yu",
	"я": "ya",
	" ": "_",
}

func isOneOf(s string, args ...string) bool {
	s = strings.ToLower(s)
	if s[0] != '/' {
		return false
	}
	// check if s has prefix args[all]
	for _, arg := range args {
		if strings.HasPrefix(s, arg) || strings.HasPrefix(s, "/"+arg) { // to write isOneOf(s, "start") instead of isOneOf(s, "/start")
			return true
		}
	}
	return false
}

func voiceTranslit(s string) string {
	// transliterates cyrillic voice names to latin
	s = strings.ToLower(s)
	// дима / димон / дим / дмитрий to dmitry
	s = strings.ReplaceAll(s, "дима", "dmitry")
	s = strings.ReplaceAll(s, "димон", "dmitry")
	s = strings.ReplaceAll(s, "дим", "dmitry")
	s = strings.ReplaceAll(s, "дмитрий", "dmitry")
	// света / светлана to svetlana
	s = strings.ReplaceAll(s, "света", "svetlana")
	s = strings.ReplaceAll(s, "светлана", "svetlana")
	// now cyrillic to latin transliteration
	for cyrillic, latin := range translitMap {
		s = strings.ReplaceAll(s, cyrillic, latin)
	}
	// now if there was ий, it translated to iy, we need to replace it with y
	s = strings.ReplaceAll(s, "iy", "y")
	return s
}
