package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// PokemonAbility is a ability certain pokemon can have
type PokemonAbility struct {
	Name string `json:"name"`
}

// PokemonAbilityEntry is a single entry about an ability
type PokemonAbilityEntry struct {
	Ability PokemonAbility `json:"ability"`
}

// PokemonMove is a move certain pokemon can have
type PokemonMove struct {
	Name string `json:"name"`
}

// PokemonMoveEntry an entry about a move a pokemon has
type PokemonMoveEntry struct {
	Move PokemonMove `json:"move"`
}

// PokemonType is a single type
type PokemonType struct {
	Name string `json:"name"`
}

// PokemonTypeEntry Single entry for a type
type PokemonTypeEntry struct {
	Type PokemonType `json:"type"`
}

// PokemonSprites are images associated with the pokemon
type PokemonSprites struct {
	BackFemale       string `json:"back_female"`
	BackShinyFemale  string `json:"back_shiny_female"`
	BackDefault      string `json:"back_default"`
	FrontFemale      string `json:"front_female"`
	FrontShinyFemale string `json:"front_shiny_female"`
	BackShiny        string `json:"back_shiny"`
	FrontDefault     string `json:"front_default"`
	FrontShiny       string `json:"front_shiny"`
}

// PokemonSearchResponse is a response we get from searching a certain pokemon
type PokemonSearchResponse struct {
	Abilities      []PokemonAbilityEntry `json:"abilities"`
	Moves          []PokemonMoveEntry    `json:"moves"`
	Types          []PokemonTypeEntry    `json:"types"`
	Weight         int                   `json:"weight"`
	Name           string                `json:"name"`
	Height         int                   `json:"height"`
	BaseExperience int                   `json:"base_experience"`
	Sprites        PokemonSprites        `json:"sprites"`
	Detail         string                `json:"detail"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var pokedexCommand = BotCommand{
	name:        "Pokedex",
	description: "Type /pokedex followed by the pokemon's name you're looking for. Pokemon with special characters in their name have those removed, except for any with a dash '-'. If you're looking for a form specific pokemon(including m/f nidoran), type dash followed by the specific form immediately after the name: /pokedex charizard-mega-x",
	matcher:     messageContainsCommandMatcher("pokedex"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		pokemon := getContentFromCommand(update.Message.Text, "pokedex")

		if pokemon == "" {
			return
		}

		searchURL := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(pokemon)
		resp, err := http.Get(searchURL)

		if err != nil {
			bot.errorReport.Log("Error Searching Pokedex: " + err.Error())
			respChan <- *NewTextBotResponse("Error Searching Pokedex", update.Message.Chat.ID)
			return
		}

		defer resp.Body.Close()

		r := PokemonSearchResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &r)
		if err != nil {
			bot.errorReport.Log("Error Parsing Pokedex Response: " + err.Error())
			respChan <- *NewTextBotResponse("Error Reading Response From Pokedex", update.Message.Chat.ID)
			return
		}

		if r.Detail == "Not found." {
			respChan <- *NewTextBotResponse(pokemon+" was not found in the Pokedex", update.Message.Chat.ID)
			return
		}

		var returnMsg bytes.Buffer

		returnMsg.WriteString("<b>")
		returnMsg.WriteString(strings.ToUpper(r.Name))
		returnMsg.WriteString("</b>\n<i>")

		// Get the types
		for i := 0; i < len(r.Types); i++ {
			returnMsg.WriteString(r.Types[i].Type.Name)
			if i < len(r.Types)-1 {
				returnMsg.WriteString(" - ")
			}
		}

		// Weight
		returnMsg.WriteString(" type\n</i>Weight: ")
		returnMsg.WriteString(strconv.FormatFloat(float64(r.Weight)/10.0, 'f', -1, 32))
		returnMsg.WriteString("kg\n")

		// Height
		returnMsg.WriteString("Height: ")
		returnMsg.WriteString(strconv.FormatFloat(float64(r.Height)/10.0, 'f', -1, 32))
		returnMsg.WriteString("m\n")

		// Base experience
		returnMsg.WriteString("Base Exp: ")
		returnMsg.WriteString(strconv.Itoa(r.BaseExperience))
		returnMsg.WriteString("\n")

		// Get the moves
		returnMsg.WriteString("\nMoves: <i>")
		numberMovesToList := min(len(r.Moves), 4)
		for i := 0; i < numberMovesToList; i++ {
			returnMsg.WriteString(r.Moves[i].Move.Name)
			if i < numberMovesToList-1 {
				returnMsg.WriteString(", ")
			}
		}

		if len(r.Moves) > 4 {
			returnMsg.WriteString(", and ")
			returnMsg.WriteString(strconv.Itoa(len(r.Moves) - 4))
			returnMsg.WriteString(" others")
		}

		// Get the moves
		returnMsg.WriteString("</i>\n\nAbilities: <i>")
		numberMovesToList = min(len(r.Abilities), 4)
		for i := 0; i < numberMovesToList; i++ {
			returnMsg.WriteString(r.Abilities[i].Ability.Name)
			if i < numberMovesToList-1 {
				returnMsg.WriteString(", ")
			}
		}

		if len(r.Abilities) > 4 {
			returnMsg.WriteString(", and ")
			returnMsg.WriteString(strconv.Itoa(len(r.Abilities) - 4))
			returnMsg.WriteString(" others")
		}

		returnMsg.WriteString("</i>\n\n")
		returnMsg.WriteString(r.Sprites.FrontDefault)

		respChan <- *NewTextBotResponse(returnMsg.String(), update.Message.Chat.ID)
	},
}
