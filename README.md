# Advent of Code Helper

This is a cli to scaffold a day's advent of code file and directory structure. When run, it will fetch your puzzle input for you, and scaffold the starter code needed to solve the problem.

Supported Languages: `python`, `golang`, `javascript(bun)`

Sample usage: 
- Run `go run . -day 2 -lang go`
- The following directory structure will be scaffolded:
    - `/day2`
        - `puzzle.txt` -> your puzzle input
        - `sample.txt` -> empty .txt file for the sample puzzle
        - `solution.go` -> starter go code

Setup:
- create a .env file in this project 
    - add `session="<your session token from advent of code website>"`
- if you want to specify a different directory for your solutions:
    - run in terminal `export ABS_PATH="<put the absolute path to the directory here>"`
- use the cli and it will scaffold your code :) 

Flags
- -lang -> accepts `python`, `golang`, `javascript`
- -day -> accepts any day number (only enter days that are available please)
- -overwrite -> include if you would like to overwrite the existing directory if one exists