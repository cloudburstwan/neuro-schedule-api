package parser

import (
	"fmt"
	"neuro-schedule-api/src/hub"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var scheduleRegex = regexp.MustCompile("<t:(\\d+):[FD]> - (?:<t:(\\d+):R> - )?(.+)")
var discordEmojiRegex = regexp.MustCompile("<a?:.+:\\d+>")

// Parse will attempt to parse a string into a hub.NeuroScheduleEntry.
// It can create a new schedule or modify an existing one, with the requirement for creating a new schedule being
// that there are 5 or more entries that match the schedule regex.
// If no matches are found, it will attempt to call the OpenAI API to attempt a natural language interpretation.
func Parse(mc string, prevSchedule []hub.NeuroScheduleEntry) ([]hub.NeuroScheduleEntry, error) {
	var result []hub.NeuroScheduleEntry
	mode := "new"

	matches := scheduleRegex.FindAllSubmatch([]byte(mc), -1)
	fmt.Println(len(matches))
	if len(matches) < 6 {
		result = prevSchedule
		mode = "update"
	}
	if len(matches) == 0 {
		var err error
		// FIXME: The natural language interpretation could be buggy due to the nature of LLMs. This will most likely need to be solved at some point.
		result, err = getScheduleModificationFromAmbiguousScheduleUpdate(result, mc)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	for _, match := range matches {
		unixSeconds, err := strconv.ParseInt(string(match[1]), 10, 64)
		if err != nil {
			return nil, err
		}

		startTime := time.Unix(unixSeconds, 0)
		isLive := len(match[2]) > 0
		var title string
		var streamers []string
		if isLive {
			title = strings.Trim(stripNonValidCharacters(string(match[3])), " ")
			streamers, err = parseStreamers(title)
			if err != nil {
				return nil, err
			}
		} else {
			title = ""
			streamers = nil
		}

		if mode == "update" {
			startDate := startTime.Format(time.DateOnly)
			existingEntryIndex := slices.IndexFunc(result, func(s hub.NeuroScheduleEntry) bool {
				return strings.HasPrefix(s.Timestamp, startDate)
			})
			if existingEntryIndex != -1 {
				result[existingEntryIndex] = hub.NeuroScheduleEntry{
					Title:     title,
					Timestamp: startTime.Format(time.RFC3339),
					Live:      isLive,
					Streamers: streamers,
				}
			} else {
				result = append(result, hub.NeuroScheduleEntry{
					Title:     title,
					Timestamp: startTime.Format(time.RFC3339),
					Live:      isLive,
					Streamers: streamers,
				})
			}
		} else {
			result = append(result, hub.NeuroScheduleEntry{
				Title:     title,
				Timestamp: startTime.Format(time.RFC3339),
				Live:      isLive,
				Streamers: streamers,
			})
		}
	}

	return result, nil
}

type streamerTitleMatchesStruct struct {
	neuro []string
	evil  []string
	vedal []string
}

var streamerTitleMatches = streamerTitleMatchesStruct{
	neuro: []string{"neuro stream", "neuro karaoke", "neuro plays"},
	evil:  []string{"evil stream", "evil karaoke", "evil plays"},
	vedal: []string{"dev stream"},
}

// parseStreamers attempts to get the list of streamers. First it tries basic formats defined in `streamerTitleMatches`.
// If that fails, it will make a call to the OpenAI API for LLM processing.
func parseStreamers(title string) ([]string, error) {
	if slices.ContainsFunc(streamerTitleMatches.neuro, func(s string) bool {
		return strings.HasPrefix(title, s)
	}) {
		return []string{"Neuro"}, nil
	}
	if slices.ContainsFunc(streamerTitleMatches.evil, func(s string) bool {
		return strings.HasPrefix(title, s)
	}) {
		return []string{"Evil"}, nil
	}
	if slices.ContainsFunc(streamerTitleMatches.vedal, func(s string) bool {
		return strings.HasPrefix(title, s)
	}) {
		return []string{"Vedal"}, nil
	}

	result, err := getStreamersFromAmbiguousStreamName(title)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var validCharactersString = "abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()-=_+'\"|{}[]?/.,<>£:;`~ "
var validCharactersArray = strings.Split(validCharactersString, "")

// stripNonValidCharacters will remove characters from the input string that are not defined in `validCharacterString`
// This is useful for e.g. stripping Unicode emojis from stream titles.
func stripNonValidCharacters(s string) string {
	var results []string
	for _, char := range discordEmojiRegex.ReplaceAllString(s, "") {
		if slices.Contains(validCharactersArray, strings.ToLower(string(char))) {
			results = append(results, string(char))
		}
	}
	return strings.Join(results, "")
}
