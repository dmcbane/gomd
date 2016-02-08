package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nochso/gomd/eol"
	"github.com/toqueteos/webbrowser"
	"gopkg.in/alecthomas/kingpin.v2"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type InputArgs struct {
	Port *int
	File *string
}

var args = InputArgs{
	Port: kingpin.Flag("port", "Listening port used by webserver").Short('p').Default("1110").Int(),
	File: kingpin.Arg("file", "Markdown file").String(),
}

type EditorView struct {
	File              string
	Content           string
	LineEndings       map[int]string
	CurrentLineEnding eol.LineEnding
}

func NewEditorView(filepath string, content string) *EditorView {
	return &EditorView{
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
	staticHTTPBox := rice.MustFindBox("static").HTTPBox()
	staticServer := http.StripPrefix("/static/", http.FileServer(staticHTTPBox))

	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("base").Parse(templateBox.MustString("base.html"))),
	}
	e.SetRenderer(t)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Get("/static/*", func(c *echo.Context) error {
		staticServer.ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})

	edit := e.Group("/edit")
	edit.Get("/*", EditHandler)
	edit.Post("/*", EditHandlerPost)

	go WaitForServer()
	e.Run(fmt.Sprintf("127.0.0.1:%d", *args.Port))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func EditHandler(c *echo.Context) error {
	var ev *EditorView
	ev, ok := c.Get("editorView").(*EditorView)
	if !ok {
		log.Println("reading file")
		filepath := c.P(0)
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to read requested file")
		}
		ev = NewEditorView(filepath, string(content))
		ev.CurrentLineEnding = eol.DetectDefault(ev.Content, eol.OSDefault())
		log.Println(ev.CurrentLineEnding.Description())
	}
	return c.Render(http.StatusOK, "base", ev)
}

func EditHandlerPost(c *echo.Context) error {
	filepath := c.P(0)
	c.Request().ParseForm()
	form := c.Request().PostForm
	eolIndex, _ := strconv.Atoi(form.Get("eol"))
	content := form.Get("content")
	convertedContent, err := eol.LineEnding(eolIndex).Apply(content)
	if err != nil {
		convertedContent = content
		log.Println("Error while converting EOL. Saving without conversion.")
	}
	ioutil.WriteFile(filepath, []byte(convertedContent), 0644)
	c.Set("editorView", NewEditorView(filepath, content))
	return EditHandler(c)
}

func WaitForServer() {
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
