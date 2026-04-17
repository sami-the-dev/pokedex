# Pokedex CLI

A command-line tool written in Go that acts as a virtual Pokédex. This application interacts with the [PokéAPI](https://pokeapi.co/) to provide information about the Pokémon world directly in your terminal.

## Features

- **Interactive REPL Environment**: Launch the CLI to a custom `Pokedex > ` prompt that continuously accepts commands.
- **API Caching**: Includes a native, in-memory caching mechanism (`pokecache`) that caches PokéAPI responses. This ensures blazing fast data lookups on subsequent identical queries.
- **Dynamic Commands**: Easily explore locations and catch Pokémon.

## Available Commands

* `help` - Show the help message and list of available commands.
* `exit` - Exit the Pokedex safely.
* `map` - Paginate forward and display the names of the next 20 location areas in the Pokémon world.
* `explore <location_area>` - Explore a specific location area to list all the Pokémon that can be encountered there.
* `catch <pokemon_name>` - Attempt to catch a specific Pokémon and display information about it (e.g., its available moves).

## Getting Started

### Prerequisites

You must have [Go](https://go.dev/) installed on your machine.

### Installation & Running

1. Clone or download the repository to your local machine.
2. Navigate to the root directory of the project in your terminal.
3. Run the CLI directly using:

```bash
# Run without building
go run main.go
```

Alternatively, you can build the binary:

```bash
# Compile to an executable
go build -o pokedexcli

# Run the executable
./pokedexcli
```

## Technologies Used

- **Go (Golang)**: Core language handling networking, REPL parsing, and object representations.
- **PokéAPI**: A full RESTful API linking to vast amounts of Pokémon data.   

## Design Decisions

- **Modularity:** Caching logic is kept strictly partitioned into the `pokecache` package avoiding business logic mixture inside the primary prompt runner.
- **Nested Struct Projection:** API responses map seamlessly onto highly representative recursive structs, abstracting the raw JSON payloads into easy-to-use Go objects.
