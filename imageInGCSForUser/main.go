package storage

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	storageLog "google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"html/template"
	"io"
	"log"
	"net/http"
)

const goBucket = "gotraining-1271.appspot.com"

func init() {
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/show", showHandler)
}

func userHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		file, header, err := req.FormFile("image")
		logError(err)
		userName := req.FormValue("userName")
		saveFile(req, userName, header.Filename, file)
		http.Redirect(res, req, "/show?userName="+userName, http.StatusFound)
		return
	}

	
	tpl := template.Must(template.ParseFiles("user.html"))
	err := tpl.Execute(res, nil)
	logError(err)
}

func saveFile(req *http.Request, userName string, fileName string, file io.Reader) {
	fileName = userName + "/" + fileName

	
	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	logStorageError(ctx, "Could not create a new client", err)
	defer client.Close()

	writer := client.Bucket(goBucket).Object(fileName).NewWriter(ctx)
	writer.ACL = []storage.ACLRule{{
		storage.AllUsers,
		storage.RoleReader}}

	io.Copy(writer, file)
	writer.Close()
}

func showHandler(res http.ResponseWriter, req *http.Request) {

	// Creating new context and client.
	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	logStorageError(ctx, "Could not create a new client", err)
	defer client.Close()

	//Parsing the template
	tpl := template.Must(template.ParseFiles("index.html"))
	err = tpl.Execute(res, getPhotoNames(ctx, client, getUserName(req)))
	logError(err)
}

func getUserName(req *http.Request) string {
	return req.FormValue("userName")
}
func getPhotoNames(ctx context.Context, client *storage.Client, userName string) []string {

	query := &storage.Query{
		Delimiter: "/",
		Prefix:    userName + "/",
	}
	objs, err := client.Bucket(goBucket).List(ctx, query)
	logError(err)

	var names []string
	for _, result := range objs.Results {
		names = append(names, result.Name)
	}
	return names
}
func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}
func logStorageError(ctx context.Context, errMessage string, err error) {
	if err != nil {
		storageLog.Errorf(ctx, errMessage, err)
	}
}
