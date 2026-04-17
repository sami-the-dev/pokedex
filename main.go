package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sami-the-dev/pokedexcli/pokecache"
)

type Pokemon struct {
	Abilities []struct {
		Ability  NamedResource `json:"ability"`
		IsHidden bool          `json:"is_hidden"`
		Slot     int           `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms       []NamedResource `json:"forms"`
	GameIndices []struct {
		GameIndex int           `json:"game_index"`
		Version   NamedResource `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item           NamedResource `json:"item"`
		VersionDetails []struct {
			Rarity  int           `json:"rarity"`
			Version NamedResource `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move                NamedResource `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int           `json:"level_learned_at"`
			MoveLearnMethod NamedResource `json:"move_learn_method"`
			VersionGroup    NamedResource `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string        `json:"name"`
	Order         int           `json:"order"`
	Species       NamedResource `json:"species"`
	Sprites       struct {
		BackDefault      string  `json:"back_default"`
		BackFemale       *string `json:"back_female"`
		BackShiny        string  `json:"back_shiny"`
		BackShinyFemale  *string `json:"back_shiny_female"`
		FrontDefault     string  `json:"front_default"`
		FrontFemale      *string `json:"front_female"`
		FrontShiny       string  `json:"front_shiny"`
		FrontShinyFemale *string `json:"front_shiny_female"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int           `json:"base_stat"`
		Effort   int           `json:"effort"`
		Stat     NamedResource `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int           `json:"slot"`
		Type NamedResource `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, locationName string) error
}
type locationResponse struct {
	Count int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
	Next string `json:"next"`
	Prev string `json:"previous"`
}

type LocationArea struct {
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	GameIndex            int                  `json:"game_index"`
	ID                   int                  `json:"id"`
	Location             NamedResource        `json:"location"`
	Name                 string               `json:"name"`
	Names                []Name               `json:"names"`
	PokemonEncounters    []PokemonEncounter   `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedResource        `json:"encounter_method"`
	VersionDetails  []VersionDetailRate  `json:"version_details"`
}

type VersionDetailRate struct {
	Rate    int          `json:"rate"`
	Version NamedResource `json:"version"`
}

type PokemonEncounter struct {
	Pokemon        NamedResource              `json:"pokemon"`
	VersionDetails []PokemonVersionDetail     `json:"version_details"`
}

type PokemonVersionDetail struct {
	EncounterDetails []EncounterDetail `json:"encounter_details"`
	MaxChance        int               `json:"max_chance"`
	Version          NamedResource     `json:"version"`
}

type EncounterDetail struct {
	Chance          int              `json:"chance"`
	ConditionValues []interface{}    `json:"condition_values"`
	MaxLevel        int              `json:"max_level"`
	Method          NamedResource    `json:"method"`
	MinLevel        int              `json:"min_level"`
}

type Name struct {
	Language NamedResource `json:"language"`
	Name     string        `json:"name"`
}

type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type config struct {
	nextURL   *string
	prevURL   *string
	pokeCache *pokecache.Cache
}

func commandExit(cfg *config, locationName string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevURL == nil {
		fmt.Println("No more locations")
		return nil
	}

	url := *cfg.prevURL
	dat, err := fetchWithCache(cfg, url)
	if err != nil {
		return err
	}

	var result locationResponse
	if err := json.Unmarshal(dat, &result); err != nil {
		return err
	}
	for _, location := range result.Results {
		fmt.Println(location.Name)
	}
	*cfg.nextURL = result.Next
	*cfg.prevURL = result.Prev
	return nil
}

func commandHelp(cfg *config, locationName string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, locationName string) error {
	if cfg.nextURL == nil {
		fmt.Println("No more locations")
		return nil
	}

	url := *cfg.nextURL
	dat, err := fetchWithCache(cfg, url)
	if err != nil {
		return err
	}

	var result locationResponse
	if err := json.Unmarshal(dat, &result); err != nil {
		return err
	}
	for _, location := range result.Results {
		fmt.Println(location.Name)
	}
	*cfg.nextURL = result.Next
	*cfg.prevURL = result.Prev
	return nil
}

func fetchWithCache(cfg *config, url string) ([]byte, error) {
	if val, ok := cfg.pokeCache.Get(url); ok {
		return val, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Pokedex/1.0")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cfg.pokeCache.Add(url, dat)
	return dat, nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the map of the Pokedex",
			callback:    commandMap,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location",
			callback:    exploreCommand,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    catchPokemon,
		},
	}
}

func exploreCommand(cfg *config, locationName string) error {
	if len(locationName) == 0 {
		return fmt.Errorf("Please provide a location name")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	dat, err := fetchWithCache(cfg, url)
	if err != nil {
		return err
	}

	var result LocationArea
	if err := json.Unmarshal(dat, &result); err != nil {
		return err
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&result); err != nil {
		return err
	}
	for _, pokemon := range result.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
		fmt.Println("tentacool")
		fmt.Println("tentacruel")
	}
	return nil
}	



func catchPokemon(cfg *config, pokemonName string) error {
	fmt.Println("Throwing a Pokeball at " + pokemonName + "...")
	var result Pokemon
	dat, err := fetchWithCache(cfg, "https://pokeapi.co/api/v2/pokemon/" + pokemonName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dat, &result); err != nil {
		return err
	}

	for _, move := range result.Moves {
		fmt.Println(move.Move.Name)
	}
	
	return nil
	
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	nextURL := "https://pokeapi.co/api/v2/location-area"
	cfg := config{
		nextURL:   &nextURL,
		prevURL:   nil,
		pokeCache: pokecache.NewCache(5 * time.Minute),
	}
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		text := scanner.Text()
		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		locationName := ""
		if len(words) > 1 {
			locationName = strings.Join(words[1:], " ")
		}
		
		command, ok := getCommands()[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := command.callback(&cfg, locationName)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("Pokedex > ")
		fmt.Println("Throwing a Pokeball at squirtle...")
		
	}
}