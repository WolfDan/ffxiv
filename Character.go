package ffxiv

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aerogo/http/client"
)

// Character represents a Final Fantasy XIV character.
type Character struct {
	Nick      string
	Server    string
	Class     string
	Level     int
	ItemLevel int
}

// GetCharacter fetches character data for a given character ID.
func GetCharacter(id string) (*Character, error) {

	url := fmt.Sprintf("https://na.finalfantasyxiv.com/lodestone/character/%s", id)
	response, err := client.Get(url).End()

	if err != nil {
		return nil, err
	}

	page := response.Bytes()
	reader := bytes.NewReader(page)
	document, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	characterName := document.Find(".frame__chara__name").Text()

	if characterName == "" {
		return nil, errors.New("Error parsing character name")
	}

	characterServer := document.Find(".frame__chara__world").Text()

	if characterServer == "" {
		return nil, errors.New("Error parsing character server")
	}

	characterLevel := document.Find(".character__class__data").Text()

	if characterLevel == "" {
		return nil, errors.New("Error parsing character level")
	}

	characterWeapon := document.Find(".db-tooltip__item__category").Text()

	if characterWeapon == "" {
		return nil, errors.New("Error parsing character class")
	}

	className := strings.Split(characterWeapon, "'")[0]
	className = strings.Replace(className, "Two-handed", "", -1)
	className = strings.Replace(className, "One-handed", "", -1)
	className = strings.Trim(className, "")

	isTwoHandedWeapon := strings.Contains(characterWeapon, "Two-handed")

	characterItems := document.Find(".db-tooltip .db-tooltip__wrapper .item_detail_box")

	itemsCount := 14

	var itemsIlvl [14]int

	hasSoul := false

	if characterItems.Length() > 1 {
		// data is duplicated on html, so we limit the amount of items
		counter := itemsCount
		characterItems.EachWithBreak(func(i int, item *goquery.Selection) bool {
			if counter <= 0 {
				return false
			}

			// ignore
			if item.Find(".db-tooltip__item__category").Text() == "Soul Crystal" {
				counter--
				hasSoul = true
				return true
			}

			re := regexp.MustCompile("[0-9]+")
			itemLevelText := re.FindStringSubmatch(item.Find(".db-tooltip__item__level").Text())[0]

			itemLevel, err := strconv.Atoi(itemLevelText)

			if err != nil {
				return false
			}

			itemsIlvl[counter] = itemLevel

			counter--

			return true
		})
	}

	iLvlSum := 0
	for _, num := range itemsIlvl {
		iLvlSum += num
	}

	if isTwoHandedWeapon {
		itemsCount-- // No shields
	}

	if hasSoul {
		itemsCount-- // Has soul stone
	}

	Ilvl := iLvlSum / itemsCount

	level, err := strconv.Atoi(characterLevel)

	if err != nil {
		return nil, err
	}

	character := &Character{Level: level, ItemLevel: Ilvl, Nick: characterName, Class: className, Server: characterServer}

	return character, nil
}
