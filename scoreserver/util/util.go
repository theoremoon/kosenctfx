package util

import "strings"

func DiscordString(s string) string {
	return strings.ReplaceAll(s, "`", "")
}
