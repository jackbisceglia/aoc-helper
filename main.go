package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Command struct {
	URL string
}

var YEAR = "2022"

var cookie_slice = []http.Cookie{
	{
		Name:  "_ga",
		Value: os.Getenv("_ga"),
	},
	{
		Name:  "_gid",
		Value: os.Getenv("_gid"),
	},
	{
		Name:  "session",
		Value: os.Getenv("session"),
	},
	{
		Name:  "_ga_MHSNPJKWC7",
		Value: os.Getenv("_ga2"),
	},
}

type Code struct {
	Day string
}

func grabInput(day string) string {
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

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

// dir -> /day${number}
//			-> index.py
//			-> puzzle.txt

type Lang struct {
	extension    string
	templateName string
}

func switchOnLang(lang string) Lang {

	switch lang {
	case "typescript":
		return Lang{
			extension:    "ts",
			templateName: "typescript.templ",
		}
	case "python":
		fallthrough
	default:
		return Lang{
			extension:    "py",
			templateName: "python.templ",
		}
	}
}

func main() {
	// setup flags
	day := flag.String("day", "-1", "Advent of Code Day")
	lang := flag.String("lang", "python", "Language")
	dev := flag.Bool("dev", false, "Should Dev")
	flag.Parse()

	fmt.Printf("%s\n", *day)

	language := switchOnLang(*lang)

	// create directories and files
	newDirName := fmt.Sprintf("day%s", *day)
	newCodeFileName := fmt.Sprintf("%s/solution.%s", newDirName, language.extension)
	newPuzzleFileName := fmt.Sprintf("%s/puzzle.txt", newDirName)

	// dev only logging
	if *dev {
		fmt.Printf("dir name: %s\n", newDirName)
		fmt.Printf("puzzle name: %s\n", newPuzzleFileName)
		fmt.Printf("code name: %s\n", newCodeFileName)
	}

	// if this dir already exists, tell the user
	// if in dev, we overwrite it, otherwise we kill the program
	dirAlreadyExists, err := exists(newDirName)
	if err != nil {
		log.Fatal(err)
	}
	if dirAlreadyExists {
		fmt.Println("Directory already exists for this day")
		if *dev {
			fmt.Println("Removing...")
			os.RemoveAll(newDirName)
		} else {
			fmt.Println("Please try again in a directory where this does not exist...")
			return
		}
	}

	if err := os.Mkdir(newDirName, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// grab puzzle
	input := grabInput(*day)

	// write puzzle
	puzzleFile, err := os.Create(newPuzzleFileName)
	if err != nil {
		log.Fatal(err)
	}
	puzzleFile.WriteString(input)

	// write code file file
	file, err := os.Create(newCodeFileName)
	if err != nil {
		log.Fatal(err)
	}

	codeTemplate := Code{
		Day: *day,
	}

	var templateFile = language.templateFile
	tmpl, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(templateFile, codeTemplate)
	if err != nil {
		panic(err)
	}

}
