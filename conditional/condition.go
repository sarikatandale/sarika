package main

import (
	"log"
	"os"
	"text/template"
)

type state struct {
	Name string
	Population  int
}

type country struct {
	state
	isCountry bool
}

func main() {
	c1 := country{
		state: state{
			Name: "Maharashtra",
			Population:  2309059,
		},
		isCountry: false,
	}

	tpl, err := template.ParseFiles("tpl.html")
	if err != nil {
		log.Fatalln(err)
	}

	err = tpl.Execute(os.Stdout, p1)
	if err != nil {
		log.Fatalln(err)
	}

}
