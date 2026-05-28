package http

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"
)

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

func (s *Server) handleGetScheduleRSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	rssItems := make([]RSSItem, 0)
	for _, stream := range s.hub.GetSchedule() {
		if stream.Live == false {
			continue
		}
		guid := strings.Split(stream.Timestamp, "T")[0]
		timestamp, err := time.Parse(time.RFC3339, stream.Timestamp)
		if err != nil {
			continue
		}

		rssItems = append(rssItems, RSSItem{
			Title:       stream.Title,
			Link:        `https://twitch.tv/vedal987`,
			Description: stream.Title,
			GUID:        guid,
			PubDate:     timestamp.Format(time.RFC1123Z),
		})
	}

	feed := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:         "Neuro-sama Schedule [Unofficial RSS feed]",
			Link:          `https://twitch.tv/vedal987`,
			Description:   "Unofficial RSS feed for Neuro-sama's schedule. Information available in the \"Schedule API\" project post in Neuro-sama Headquarters",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         rssItems,
		},
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		http.Error(w, "Internal Server Error while packaging response", http.StatusInternalServerError)
	}

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(feed); err != nil {
		http.Error(w, "Internal Server Error while encoding RSS", http.StatusInternalServerError)
	}
	err = encoder.Close()
	if err != nil {
		return
	}
}
