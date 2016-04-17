package storage

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	storageLog "google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const cloudBucket = "golangcourse.appspot.com"
const textFile = "textFile.txt"

func init() {
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {

	
	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	logStorageError(ctx, "Could not create a new client", err)
	defer client.Close()
	store(ctx, client, " folder/")
	store(ctx, client, " folderlike+")
	printFolders(ctx, client, res, "/")
	printFolders(ctx, client, res, "+")
}

func printFolders(ctx context.Context, client *storage.Client, res http.ResponseWriter, delimeter string) {
	fmt.Fprintf(res, "Delimeter ["+delimeter+"]\n")
	query := &storage.Query{
		Delimiter: delimeter,
	}
	objs, err := client.Bucket(cloudBucket).List(ctx, query)
	logError(err)
	for _, subfolder := range objs.Prefixes {
		fmt.Fprintf(res, "Folder: "+subfolder+"\n")
	}
}


func store(ctx context.Context, client *storage.Client, folderPostfix string) {
	
	reader, err := os.Open(textFile)
	logError(err)

		for i := 0; i < 3; i++ {
		writer := client.Bucket(cloudBucket).Object(strconv.Itoa(i) + folderPostfix + textFile).NewWriter(ctx)
		writer.ACL = []storage.ACLRule{{
			storage.AllUsers,
			storage.RoleReader}}
		io.Copy(writer, reader)
		writer.Close()
	}
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
