package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/feeds"
)

const tagesschauApiUrl = "https://www.tagesschau.de/api2u/news"
const iconUrl = "https://upload.wikimedia.org/wikipedia/commons/thumb/0/06/Tagesschau.de_logo.svg/640px-Tagesschau.de_logo.svg.png"
const port = ":5050"

func getEilMeldungen() ([]News, error) {
	var news []News

	response, err := http.Get(tagesschauApiUrl)
	if err != nil {
		return news, err
	}

	var newsResponse NewsResponse
	if err = json.NewDecoder(response.Body).Decode(&newsResponse); err != nil {
		return news, err
	}

	for _, entry := range newsResponse.News {
		for _, tag := range entry.Tags {
			if tag.Tag == "Eilmeldung" {
				news = append(news, entry)
				break
			}
		}
	}

	return news, nil
}

func getFeed(w http.ResponseWriter) (feeds.Feed, error) {
	feed := feeds.Feed{}

	feed.Title = "Tagesschau Eilmeldungen"
	feed.Description = "Wichtige Nachrichten der Tagesschau"
	feed.Link = &feeds.Link{Href: "https://www.tagesschau.de"}
	feed.Image = &feeds.Image{Title: "Tagesschau", Url: iconUrl, Link: iconUrl}

	news, err := getEilMeldungen()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return feed, err
	}

	for _, entry := range news {
		description := fmt.Sprintf("<![CDATA[<img src=\"%s\"/> %s", entry.TeaserImage.ImageVariants.One6X9256, entry.FirstSentence)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: entry.ShareURL},
			Id:          entry.ShareURL,
			Description: description,
			Created:     entry.Date,
		})
	}

	return feed, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service up.")
	})

	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		if feed, err := getFeed(w); err == nil {
			rss, _ := feed.ToRss()
			w.Write([]byte(rss))
		}
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		if feed, err := getFeed(w); err == nil {
			json, _ := feed.ToJSON()
			w.Write([]byte(json))
		}
	})

	http.HandleFunc("/atom", func(w http.ResponseWriter, r *http.Request) {
		if feed, err := getFeed(w); err == nil {
			atom, _ := feed.ToAtom()
			w.Write([]byte(atom))
		}
	})

	fmt.Printf("Listening on %s\n", port)
	http.ListenAndServe(port, nil)
}
