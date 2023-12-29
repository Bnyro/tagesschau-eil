package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/feeds"
)

const tagesschauApiUrl = "https://www.tagesschau.de/api2u/news"
const port = ":5050"
const minAmountOfEntries = 20
const maxAmountOfEntries = 100
const eilMeldungTag = "Eilmeldung"

var storedEilNews []News

func getEilMeldungen(apiUrl string) ([]News, string, error) {
	var news []News

	response, err := http.Get(apiUrl)
	if err != nil {
		return news, "", err
	}

	var newsResponse NewsResponse
	if err = json.NewDecoder(response.Body).Decode(&newsResponse); err != nil {
		return news, "", err
	}

	for _, entry := range newsResponse.News {
		for _, tag := range entry.Tags {
			if tag.Tag == eilMeldungTag {
				news = append(news, entry)
				break
			}
		}
	}

	return news, newsResponse.NextPage, nil
}

func appendNews(news []News) {
out:
	for _, entry := range news {
		for _, existingEntry := range storedEilNews {
			if entry.ExternalID == existingEntry.ExternalID {
				continue out
			}
		}

		storedEilNews = append(storedEilNews, entry)
	}

	sort.Slice(storedEilNews, func(i, j int) bool {
		return storedEilNews[i].Date.After(storedEilNews[j].Date)
	})

	if len(storedEilNews) > maxAmountOfEntries {
		storedEilNews = storedEilNews[len(storedEilNews)-maxAmountOfEntries:]
	}
}

func updateFeed() {
	newEntries, nextPage, err := getEilMeldungen(tagesschauApiUrl)

	if err != nil {
		return
	}

	appendNews(newEntries)

	for len(storedEilNews) < minAmountOfEntries {
		newEntries, nextPage, err = getEilMeldungen(nextPage)

		if err != nil {
			break
		}

		appendNews(newEntries)
	}
}

func getFeed(r *http.Request, w http.ResponseWriter) feeds.Feed {
	feed := feeds.Feed{}

	iconUrl := "https://" + r.Host + "/static/tagesschau.png"

	feed.Title = "Tagesschau Eilmeldungen"
	feed.Description = "Wichtige Nachrichten der Tagesschau"
	feed.Link = &feeds.Link{Href: "https://www.tagesschau.de"}
	feed.Image = &feeds.Image{Title: "Tagesschau", Url: iconUrl, Link: iconUrl}
	feed.Copyright = "ARD-aktuell / tagesschau.de"

	for _, entry := range storedEilNews {
		content := fmt.Sprintf("<img src=\"%s\" />\n%s", entry.TeaserImage.ImageVariants.Land960, entry.FirstSentence)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: entry.ShareURL},
			Id:          entry.ExternalID,
			Description: entry.FirstSentence,
			Content:     content,
			Created:     entry.Date,
		})
	}

	return feed
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service up.")
	})
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		feed := getFeed(r, w)
		rss, _ := feed.ToRss()
		w.Write([]byte(rss))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		feed := getFeed(r, w)
		json, _ := feed.ToJSON()
		w.Write([]byte(json))

	})

	http.HandleFunc("/atom", func(w http.ResponseWriter, r *http.Request) {
		feed := getFeed(r, w)
		atom, _ := feed.ToAtom()
		w.Write([]byte(atom))
	})

	go func() {
		updateFeed()

		time.Sleep(time.Second * 30)

		updateFeed()

		for range time.Tick(time.Minute * 15) {
			updateFeed()
		}
	}()

	fmt.Printf("Listening on %s\n", port)
	http.ListenAndServe(port, nil)
}
