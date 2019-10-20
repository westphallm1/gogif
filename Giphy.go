package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type GiphyResponse struct {
	Data []GifResponse
}

type GifResponse struct {
	Id    string
	Title string
}

const URL_FMT = "http://api.giphy.com/v1/gifs/search?api_key=%s&q=%s"

// TODO this is flimsy/liable to break, but they don't expose the image directly
const GIF_FMT = "http://i.giphy.com/media/%s/100.gif"

func getGiphyJSON(query string) []GifResponse {
	safeQuery := url.QueryEscape(query)
	api_key := os.Getenv("GIPHY_KEY")
	url := fmt.Sprintf(URL_FMT, api_key, safeQuery)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var giphyResponse GiphyResponse
	err = json.Unmarshal(body, &giphyResponse)
	if err != nil {
		log.Fatal(err)
	}
	return giphyResponse.Data
}

func downloadGiphy(id string) io.Reader {
	url := fmt.Sprintf(GIF_FMT, id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body
}
