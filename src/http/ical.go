package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

var namespace uuid.UUID = uuid.MustParse("8f85716f-1c9e-4976-b830-56bae9644f38")

// streamer name -> sha1 -> uuid -> urn:uuid URI
func streamerURI(name string) string {
	return fmt.Sprintf("urn:uuid:%s", uuid.NewSHA1(namespace, []byte(name)).String())
}

func (s *Server) handleGetScheduleICS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	cal := ics.NewCalendarFor("neuro-schedule-api")
	cal.SetMethod(ics.MethodPublish)
	cal.SetName("Neuro-sama Schedule [Unofficial]")
	cal.SetDescription("Unofficial iCal feed for Neuro-sama's schedule. Information available in the \"Schedule API\" project post in Neuro-sama Headquarters")

	// TODO: if https://github.com/arran4/golang-ical/pull/144 ever gets merged then can just use that instead of this
	cal.CalendarProperties = append(cal.CalendarProperties, ics.CalendarProperty{
		BaseProperty: ics.BaseProperty{
			IANAToken: "X-APPLE-CALENDAR-COLOR",
			Value:     "#FF8AB6",
		},
	})

	for _, stream := range s.hub.GetSchedule() {
		if stream.Live == false {
			continue
		}

		guid := strings.Split(stream.Timestamp, "T")[0]
		timestamp, err := time.Parse(time.RFC3339, stream.Timestamp)
		if err != nil {
			continue
		}

		event := cal.AddEvent(guid)
		event.SetSummary(stream.Title)
		event.SetStartAt(timestamp)
		event.SetDuration(3 * time.Hour)
		event.SetURL("https://twitch.tv/vedal987")
		event.SetLocation("https://twitch.tv/vedal987")
		for _, streamer := range stream.Streamers {
			event.AddProperty(
				ics.ComponentPropertyAttendee,
				streamerURI(streamer),
				ics.WithCN(streamer),
				&ics.KeyValues{Key: "PARTSTAT", Value: []string{"ACCEPTED"}},
			)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cal.Serialize(ics.WithNewLineWindows))) // to comply with RFC5545 it has to be CRLF apparently
}
