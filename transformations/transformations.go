package transformations

import "strings"

func SpaceToUnderScore(text string) (string) {
	// take the text
	// convert it into lowercase
	// replace the spaces into _
	text = 	strings.TrimSpace(text)
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "_")

	return text
}

