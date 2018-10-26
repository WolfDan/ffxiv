package ffxiv_test

import (
	"testing"

	"./"

	"github.com/stretchr/testify/assert"
)

func TestGetCharacter(t *testing.T) {
	character, err := ffxiv.GetCharacter("9015414")

	characterResult := ffxiv.Character{
		Level:     70,
		ItemLevel: 382,
		Nick:      "Aky Otara",
		Class:     "Conjurer",
		Server:    "Asura",
	}

	assert.NoError(t, err)
	assert.Equal(t, characterResult, character)
}

// func TestGetCharacterIDFail(t *testing.T) {
// 	id, err := ffxiv.GetCharacterID("asdfasdfghjklasdf", "Asura")

// 	assert.Error(t, err)
// 	assert.Empty(t, id)
// }
