package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func CustomFilter(suggestions []prompt.Suggest, sub string) []prompt.Suggest {
	return filterSuggestions(suggestions, sub, true, strings.Contains)
}

func filterSuggestions(suggestions []prompt.Suggest, sub string, ignoreCase bool, function func(string, string) bool) []prompt.Suggest {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggestions))
	for i := range suggestions {
		c := suggestions[i].Text + suggestions[i].Description
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if function(c, sub) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}
