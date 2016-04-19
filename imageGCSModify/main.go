package storage

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	storageLog "google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"html/template"
	"log"
	"net/http"
)

const goBucket = "golangcourse.appspot.com"

func init() {
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {

	
	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	logStorageError(ctx, "Could not create a new client", err)
	defer client.Close()

	tpl := template.Must(template.ParseFiles("index.html"))
	err = tpl.Execute(res, getPhotoNames(ctx, client))
	logError(err)
}

func getPhotoNames(ctx context.Context, client *storage.Client) []string {

	query := &storage.Query{
		MaxResults: 4, 
		Delimiter:  "/",
		Prefix:     "photos/",
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
