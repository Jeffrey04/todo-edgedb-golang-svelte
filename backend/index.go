package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID         edgedb.UUID
	Item       string
	Done_at    edgedb.OptionalDateTime
	Updated_at time.Time
	Created_at time.Time
}

func db_connect() (context.Context, *edgedb.Client) {
	ctx := context.Background()

	client, err := edgedb.CreateClient(ctx, edgedb.Options{})

	if err != nil {
		log.Fatal(err)
	}

	return ctx, client
}

func index(c *gin.Context) {
	ctx, client := db_connect()

	defer client.Close()

	var result []byte

	err := client.QueryJSON(
		ctx,
		"SELECT Todo {id, item, updated_at, created_at, done_at}",
		&result,
	)

	if err != nil {
		log.Fatal(err)
	}

	c.Data(http.StatusOK, gin.MIMEJSON, result)
}

func add(c *gin.Context) {
	ctx, client := db_connect()

	defer client.Close()

	var data Todo
	c.BindJSON(&data)

	var inserted struct{ id edgedb.UUID }
	err := client.QuerySingle(
		ctx,
		"INSERT Todo {item := <str>$0, updated_at := <datetime>$1, created_at := <datetime>$2}",
		&inserted,
		data.Item,
		data.Updated_at,
		data.Created_at)

	if err != nil {
		log.Fatal(err)
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/%s/", inserted.id))
	c.Status(http.StatusCreated)
}

func update(c *gin.Context) {
	ctx, client := db_connect()

	defer client.Close()

	var data Todo
	c.BindJSON(&data)
	uuid, err := edgedb.ParseUUID(strings.Trim(c.Param("id"), "/"))

	if err != nil {
		log.Fatal(err)
	}

	var updated struct{ id edgedb.UUID }
	done_at, ok := data.Done_at.Get()

	if ok {
		err = client.QuerySingle(
			ctx,
			"UPDATE Todo FILTER .id = <uuid>$0 SET {item := <str>$1, updated_at := <datetime>$2, created_at := <datetime>$3, done_at := <datetime>$4} ",
			&updated,
			uuid,
			data.Item,
			data.Updated_at,
			data.Created_at,
			done_at)
	} else {
		err = client.QuerySingle(
			ctx,
			"UPDATE Todo FILTER .id = <uuid>$0 SET {item := <str>$1, updated_at := <datetime>$2, created_at := <datetime>$3, done_at := {}} ",
			&updated,
			uuid,
			data.Item,
			data.Updated_at,
			data.Created_at)
	}

	if err != nil {
		log.Fatal(err)
	}

	//c.Status(http.StatusNoContent)
}

func delete(c *gin.Context) {
	ctx, client := db_connect()

	defer client.Close()

	uuid, err := edgedb.ParseUUID(strings.Trim(c.Param("id"), "/"))

	if err != nil {
		log.Fatal(err)
	}

	var deleted struct{ id edgedb.UUID }
	err = client.QuerySingle(
		ctx,
		"DELETE Todo FILTER .id = <uuid>$0",
		&deleted,
		uuid)

	if err != nil {
		log.Fatal(err)
	}

	c.Status(http.StatusNoContent)
}

func main() {
	r := gin.Default()

	r.StaticFile("/", "./frontend/dist/index.html")
	r.StaticFile("/index.html", "./frontend/dist/index.html")
	r.Static("/assets", "./frontend/dist/assets")

	r.GET("/api/", index)
	r.PUT("/api/:id/", update)
	r.DELETE("/api/:id/", delete)
	r.POST("/api/", add)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
