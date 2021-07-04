package main

import (
	"encoding/json"
	"fmt"
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
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		if err := c.Render(http.StatusOK, "index.html", nil); err != nil {
			log.Println(err)
			return err
		}

		return nil
	})

	e.POST("/result", func(c echo.Context) error {
		target := c.FormValue("target")
		if len(target) == 0 {
			// [TODO]: 文字列長さが0のときの処理
			return nil
		}

		words, err := fetchWords(target)
		if err != nil {
			log.Println(err)
			return err
		}

		if err := c.Render(http.StatusOK, "result.html", words); err != nil {
			log.Println(err)
			return err
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8000"))
}

func fetchWords(target string) ([]Word, error) {
	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	// [TODO]: URLの設定方法
	url := fmt.Sprintf("http://172.30.0.3:8080/v1/roman?target=%v", target)
	req, err := http.NewRequest("GET", url, nil)
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
