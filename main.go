package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Word struct {
	Raw    string `json:"raw"`
	Roman  string `json:"roman"`
	Vowels string `json:"vowels"`
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("template/index.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		words, err := fetchWords()
		if err != nil {
			return err
		}

		if err := c.Render(http.StatusOK, "index.html", words); err != nil {
			return err
		}

		return nil
	})

	e.POST("/result", func(c echo.Context) error {
		log.Println(c.FormValue("target"))
		// renderの処理
		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}

func fetchWords() ([]Word, error) {
	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	// [TODO]: URLの設定方法
	req, err := http.NewRequest("GET", "http://172.30.0.3:8080/v1/roman?target=chiyuki", nil)
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cliID := os.Getenv("BASIC_AUTH_CLIENT_ID")
	cliSec := os.Getenv("BASIC_AUTH_CLIENT_SECRET")
	req.SetBasicAuth(cliID, cliSec)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var words []Word
	if err := json.Unmarshal(b, &words); err != nil {
		log.Println(err)
		return nil, err
	}

	return words, nil
}
