package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/joho/godotenv/autoload"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type Command struct {
	URL string
}

var YEAR = "2022"

var cookie_slice = []http.Cookie{
	{
		Name:  "session",
		Value: os.Getenv("session"),
	},
}

type Code struct {
	Day string
}

func fetchInputWithCache(day string) string {
	// check cache
	cachedPuzzlePath := fmt.Sprintf("puzzle-cache/%s-%s.txt", YEAR, day)
	exists, err := exists(cachedPuzzlePath)

	if err != nil {
		log.Fatal(err)
	}

	if exists {
		// read file and return
		file, err := os.ReadFile(cachedPuzzlePath)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Using cached input")
		return string(file)
	}

	// init client and set url
	client := &http.Client{}
	url := fmt.Sprintf("https://adventofcode.com/%s/day/%s/input", YEAR, day)

	// init req
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// attach cookies
	for _, cookie := range cookie_slice {
		req.AddCookie(&cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	puzzle := string(body)

	// add to cache
	cachedPuzzleFile, err := os.Create(cachedPuzzlePath)
	if err != nil {
		log.Fatal(err)
	}
	cachedPuzzleFile.WriteString(puzzle)

	return puzzle
}

type Lang struct {
	name      string
	extension string
}

func get_language_details(lang string) Lang {
	switch lang {
	case "go":
		return Lang{
			name:      "go",
			extension: "go",
		}
	case "python":
		return Lang{
			name:      "python",
			extension: "py",
		}
	case "javascript":
		fallthrough
	default:
		return Lang{
			name:      "javascript",
			extension: "js",
		}
	}
}

type CLI_Flags struct {
	day  string
	lang string
	dev  bool
}

func parse_flags() CLI_Flags {
	day := flag.String("day", "-1", "Advent of Code Day")
	lang := flag.String("lang", "javascript", "Language")
	dev := flag.Bool("dev", false, "Should Dev")

	flag.Parse()

	return CLI_Flags{
		day:  *day,
		lang: *lang,
		dev:  *dev,
	}
}

func main() {
	// setup flags
	flags := parse_flags()

	fmt.Printf("%s\n", flags.day)

	language := get_language_details(flags.lang)
	basePathWithBackSlash := ""
	if os.Getenv("ABS_PATH") != "" {
		basePathWithBackSlash = os.Getenv("ABS_PATH") + "/"
	}

	// create directories and files
	newDirPath := fmt.Sprintf("%sday%s", basePathWithBackSlash, flags.day)
	newCodeFilePath := fmt.Sprintf("%s/solution.%s", newDirPath, language.extension)
	newPuzzleFilePath := fmt.Sprintf("%s/puzzle.txt", newDirPath)
	newSampleFilePath := fmt.Sprintf("%s/sample.txt", newDirPath)
	templateFileName := fmt.Sprintf("%s.tmpl", language.name)
	templateFilePath := fmt.Sprintf("language-templates/%s", templateFileName)

	// if this dir already exists, tell the user
	// if in dev, we overwrite it, otherwise we kill the program
	dirAlreadyExists, err := exists(newDirPath)
	if err != nil {
		log.Fatal(err)
	}
	if dirAlreadyExists {
		fmt.Println("Directory already exists for this day")
		if flags.dev {
			fmt.Println("Removing...")
			os.RemoveAll(newDirPath)
		} else {
			fmt.Println("Please try again in a directory where this does not exist...")
			return
		}
	}

	if err := os.Mkdir(newDirPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// grab puzzle
	input := fetchInputWithCache(flags.day)

	// write sample
	_, err = os.Create(newSampleFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// write puzzle
	puzzleFile, err := os.Create(newPuzzleFilePath)
	if err != nil {
		log.Fatal(err)
	}
	puzzleFile.WriteString(input)

	// write code file file
	codeFile, err := os.Create(newCodeFilePath)
	if err != nil {
		log.Fatal(err)
	}

	code := Code{
		Day: flags.day,
	}

	tmpl, err := template.New(templateFileName).ParseFiles(templateFilePath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(codeFile, code)
	if err != nil {
		panic(err)
	}

}
