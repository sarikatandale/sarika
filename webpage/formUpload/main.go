package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)


func uploadTextFile(res http.ResponseWriter, req *http.Request) {

	temp, err := template.ParseFiles("./uploadTextFile.html")

	logError(err)

	if req.Method == "POST" {

		inFile, _, err := req.FormFile("file")
		defer inFile.Close()
	
		logError(err)

		contents, err := ioutil.ReadAll(inFile)
	
		logError(err)

		err = temp.Execute(res, string(contents))
	} else {
		err = temp.Execute(res, false)
	}

	logError(err)

}

func main() {


	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", uploadTextFile)

	log.Println("Listening to 8080 ...")
	http.ListenAndServe(":8080", nil)
}


func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
