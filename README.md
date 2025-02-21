# pokedex-cli: a Pokédex CLI tool

## What is it?

**pokedex-cli** starts with an implementation of the guided project [Build a Pokedex in Go](https://www.boot.dev/courses/build-pokedex-cli-golang) from the [boot.dev](https://www.boot.dev) platform.

**New features** were then **added** to make the project more personal, more complete and more user-friendly, the distinction between the **tutorial** features and the **new** features is made below.

It is a **CLI tool** that uses the [PokéAPI](https://pokeapi.co/) API to fetch Pokemon-related data. The game within a command line terminal allows you to explore the Pokémon universe, catch Pokémons, and display information about them.


                                                                                  @     
                                : :                                             # #     
          @ # # @               : : :                   @ # # # #             #   #     
            @ # # # @           : : :               @ # # # # # @           #     #     
              @ # # # @         : :               # @ @ @ # # #           #       #     
                @ # @ @ #       :               # : : @ @ @ #           #         #     
                  @     : #                   # : : : : : #           #           #     
                    #   : : #                 #   : : : #           #             @     
          : :         # : : : @ # # # # # # @     : : #           #             #       
        : : :           @ #                 #   : @ #           # :           #         
          : :         #                           #               # : : : : #           
            :       #   #                 # #       #               # # : : : #         
                  @   #   #             #   # #     @                   # : : : #       
                  #   # # #             # # # #       #               # : : @ @ @ #     
                @     : # :             : # # :       @             # @ @ @ @ @ #       
                # @           #                 @ @     #         # @ @ @ @ @ #         
              # @ @ @                         @ @   @   @       # @ @ @ @ @ #           
              # : @     #     # #       #       @ @     : #       # @ @ @ #             
              #           # @     # # @                 : #         # @ @ # @           
            # # @                               # #   : : #           # # # # @         
            # @ # @                           # @ @ # : : @             # # # #         
            # @ @ #     : : : : :           #   @ #     : : @             # # #         
              # @     : : : : : : :             #         : #             @ # #         
                #   : : : : : : : : :         #           @ #             # # @         
              @   : : : : : : : : : : :               @ @ @ @ @         # # #           
              #   : : : : : : : : : : :     :       @ @ @ @ @ #         # # #           
              #   : : : : : : : : : : :   :       @ @       : #       @ # # @       :   
              #   : : : : : : : : : : :   @                 : #       # # #       : :   
              # : : : : : : : : : : : : @                 : : #     @ # # @     : : :   
          @ # # @ : : : : : : : : : : : #                 : @ #     # # @     : :       
          # @ @ # : : : : : : : : : : : #             : @ @ @ # @ # # @                 
            # @ @ # : : : : : : : : : : # :         : : : @ @ # # @                     
              # @ @ # : : : : : : : : @ # @ :     : : : : : : # @                       
                # # @ # @ : : : : @ @ @ @ # : : : : : : : : @ #         :               
                          # # @ @ @ @ @ #   # : : : : : : : #           : :             
                              # # # # #       # : : : : : @               : :           
                                                # # @ : : #               : : :         
                                                      # @ @ #             : : : :       
                                                      # @ @ #               : :         
                                                        # #

## Features from the guided tutorial
The following features were added while going through the guided tutorial:
  - **REPL loop** with prompt to interact with the user
```sh
Pokedex > _
```

  - Implementation of a **cache** to remember requests results **within** the current session, with a go routine to automatically **clear** entries older than a given time interval.
  - Implementation of a list of **commands**:
    - **exit**: exit the program
    - **help**: display the list of available commands
    - **map**: navigate and display the next 20 locations
    - **explore**: list the Pokémons in a given location
    - **catch**: throw a ball and try to catch a Pokémon
    - **inspect**: display information about a caught Pokémon
    - **pokedex**: list all the Pokémons caught

## New Features
The following features were added independently, and are specific to this project:
  - **Persistence** of the cache in between sessions, the entries are saved in a binary file and loaded on program start.
  - **Persistence** of the pokedex entries in between sessions. The pokemons caught and the position on the map are loaded on program start.
  - Add support of **special keys** in the CLI terminal:
    - **Backspace** to delete the last character
    - **Left** and **Right** arrow keys to navigate in the prompt line.
    - **Up** and **Down** arrow keys to navigate in the command history.
    - **Tab** key to auto-complete the command, location or pokemon name being written.
    - **Ctrl+C** or **Ctrl+D** for a clean exit (cache and pokedex are saved).
    - **Double Tab** to show suggestions of available commands, locations or pokemon names depending on context.
  - **Improvement** of the existing commands:
    - **catch**:
      - On success, the terminal now displays an **ASCII art** of the Pokémon being caught.
      - The **failure message** is now randomnly picked from a list of options, to add variety.
      - The user can choose to use the **--great-ball**, **--ultra-ball** or **--master-ball** options to increase the success rate of the catch.
    - **pokedex**:
      - The output now lists all the pokemons caught in properly **formated columns**, with their **#number**, and with a **???** tag for pokemons that have not been caught yet.
    - **explore**:
      - The pokemons that have been caught yet are shown in the terminal with a **New!** tag
    - **map**:
      - The locations where **New!** pokemons can be found are now diplayed **in bold** in the output.
    - **inspect**:
      - Add the **ASCII art** of the pokemon before showing the stats.
  - **Addition** of new commands:
    - **compare**:
      - Two pokemons from the pokedex are shown **side by side** in the terminal, both with their **ASCII art** and their **names** and **stats**.

  - **Refactoring** of the **code** into clean packages.
  - **Refactoring** of the **tests** to be able to simulate user input in the CLI.
  - Addition of a **--debug** mode with proper **logging**.

## Installation from sources

The source code is currently hosted on GitHub at:
[www.github.com/noch-g/pokedex-cli](https://www.github.com/noch-g/pokedex-cli)

```sh
# go install
go install github.com/noch-g/pokedex-cli@latest
```

<hr>
