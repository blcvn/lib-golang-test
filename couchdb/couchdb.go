package main

import (
	"fmt"
	"net/url"

	"github.com/zemirco/couchdb"
)

// create your own document
type dummyDocument struct {
	couchdb.Document
	Foo  string `json:"foo"`
	Beep string `json:"beep"`
}

// start
func main() {
	u, err := url.Parse("http://127.0.0.1:5984/")
	if err != nil {
		panic(err)
	}

	// create a new client
	// client, err := couchdb.NewClient(u)
	username := "admin"
	password := "password"

	client, err := couchdb.NewAuthClient(username, password, u)
	if err != nil {
		panic(err)
	}

	// get some information about your CouchDB
	info, err := client.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)

	// create a database
	dbname := "dummydb"
	if _, err = client.Create(dbname); err != nil {
		panic(err)
	}

	// use your new "dummy" database and create a document
	db := client.Use(dbname)
	doc := &dummyDocument{
		Foo:  "bar",
		Beep: "bopp",
	}
	result, err := db.Post(doc)
	if err != nil {
		panic(err)
	}

	docR := &dummyDocument{}
	// get id and current revision.
	if err := db.Get(docR, result.ID); err != nil {
		panic(err)
	}

	fmt.Printf("Result from couchdb: %+v \n", docR)

	// delete document
	if _, err = db.Delete(docR); err != nil {
		panic(err)
	}

	// and finally delete the database
	if _, err = client.Delete(dbname); err != nil {
		panic(err)
	}

}
