package main

import "strings"

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
