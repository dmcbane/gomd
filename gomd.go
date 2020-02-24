package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nochso/gomd/eol"
	"github.com/toqueteos/webbrowser"
	"gopkg.in/alecthomas/kingpin.v2"
)

type inputArgs struct {
	Port *int
	File *string
}

var args = inputArgs{
	Port: kingpin.Flag("port", "Listening port used by webserver").Short('p').Default("1110").Int(),
	File: kingpin.Arg("file", "Markdown file").String(),
}

type editorView struct {
	File              string
	Content           string
	LineEndings       map[int]string
	CurrentLineEnding eol.LineEnding
}

func newEditorView(filepath string, content string) *editorView {
	return &editorView{
		File:        filepath,
		Content:     content,
		LineEndings: eol.Descriptions,
	}
}

func main() {
	// Parse command line arguments
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// Prepare (optionally) embedded resources
	templateBox := rice.MustFindBox("template")
	// staticHTTPBox := rice.MustFindBox("static").HTTPBox()
	// staticServer := http.StripPrefix("/static/", http.FileServer(staticHTTPBox))

	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("base").Parse(templateBox.MustString("base.html"))),
	}
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static/*", "static")

	edit := e.Group("/edit")
	edit.GET("/*", editHandler)
	edit.POST("/*", editHandlerPost)

	go waitForServer()
	e.Logger.Fatal(e.Start(fmt.Sprintf("%d", *args.Port)))
}

type Template struct {
	templates *template.Template
}

// Render Function
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func editHandler(c echo.Context) error {
	var ev *editorView
	ev, ok := c.Get("editorView").(*editorView)
	if !ok {
		log.Println("reading file")
		filepath := c.ParamNames()[0]
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to read requested file")
		}
		ev = newEditorView(filepath, string(content))
		ev.CurrentLineEnding = eol.DetectDefault(ev.Content, eol.OSDefault())
		log.Println(ev.CurrentLineEnding.Description())
	}
	return c.Render(http.StatusOK, "base", ev)
}

func editHandlerPost(c echo.Context) error {
	filepath := c.ParamNames()[0]
	eolIndex, _ := strconv.Atoi(c.FormValue("eol"))
	content := c.FormValue("content")
	convertedContent, err := eol.LineEnding(eolIndex).Apply(content)
	if err != nil {
		convertedContent = content
		log.Println("Error while converting EOL. Saving without conversion.")
	}
	ioutil.WriteFile(filepath, []byte(convertedContent), 0644)
	c.Set("editorView", newEditorView(filepath, content))
	return editHandler(c)
}

func waitForServer() {
	log.Printf("Waiting for listener on port %d", *args.Port)
	url := fmt.Sprintf("http://localhost:%d/edit/%s", *args.Port, url.QueryEscape(*args.File))
	for {
		time.Sleep(time.Millisecond * 50)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		resp.Body.Close()
		break
	}
	log.Println("Opening " + url)
	if err := webbrowser.Open(url); err != nil {
		log.Printf("Possible error while opening browser: %s", err)
	}
}
