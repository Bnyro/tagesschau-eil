package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/feeds"
)

const tagesschauApiUrl = "https://www.tagesschau.de/api2u/news"
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

func getFeed(r *http.Request, w http.ResponseWriter) (feeds.Feed, error) {
	feed := feeds.Feed{}

	iconUrl := "https://" + r.Host + "/static/tagesschau.png"

	feed.Title = "Tagesschau Eilmeldungen"
	feed.Description = "Wichtige Nachrichten der Tagesschau"
	feed.Link = &feeds.Link{Href: "https://www.tagesschau.de"}
	feed.Image = &feeds.Image{Title: "Tagesschau", Url: iconUrl, Link: iconUrl}
	feed.Copyright = "ARD-aktuell / tagesschau.de"

	news, err := getEilMeldungen()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return feed, err
	}

	for _, entry := range news {
		content := fmt.Sprintf("<img src=\"%s\" />\n%s", entry.TeaserImage.ImageVariants.Land256, entry.FirstSentence)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: entry.ShareURL},
			Id:          entry.ExternalID,
			Description: entry.FirstSentence,
			Content:     content,
			Created:     entry.Date,
		})
	}

	return feed, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service up.")
	})
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		if feed, err := getFeed(r, w); err == nil {
			rss, _ := feed.ToRss()
			w.Write([]byte(rss))
		}
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		if feed, err := getFeed(r, w); err == nil {
			json, _ := feed.ToJSON()
			w.Write([]byte(json))
		}
	})

	http.HandleFunc("/atom", func(w http.ResponseWriter, r *http.Request) {
		if feed, err := getFeed(r, w); err == nil {
			atom, _ := feed.ToAtom()
			w.Write([]byte(atom))
		}
	})

	fmt.Printf("Listening on %s\n", port)
	http.ListenAndServe(port, nil)
}
