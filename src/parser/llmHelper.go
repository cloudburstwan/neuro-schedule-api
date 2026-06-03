package parser

import (
	"context"
	"encoding/json"
	"neuro-schedule-api/src/hub"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var client openai.Client

type ScheduleJson struct {
	Schedule []hub.NeuroScheduleEntry `json:"schedule"`
	String   string                   `json:"string"`
}

func ConnectToOpenAI(key string) {
	client = openai.NewClient(
		option.WithAPIKey(key),
	)
}

func getStreamersFromAmbiguousStreamName(name string) ([]string, error) {
	ctx := context.Background()
	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage("Your job is to take an input string from the user and figure out who is present on the stream that the user is describing. The user's input will be short and consise. You must respond with nothing but the names of the people you have identified separated by a `,` character. Do not include spaces unless they're part of the name. You'll commonly see `Neuro`, `Evil`, and `Vedal`, however you may encounter other names that you'll need to correctly extract. If you encounter \"The Twins\" or similar variations, this is usually both Neuro and Evil and must be put into the output individually. If you encounter \"Dev Stream\" (or variations) without \"Vedal\" elsewhere in the title, add \"Vedal\" to the output regardless."),
			openai.UserMessage(name),
		},
		Model: openai.ChatModelGPT5_4Nano,
	})
	if err != nil {
		return nil, err
	}

	return strings.Split(res.Choices[0].Message.Content, ","), nil
}

func getScheduleModificationFromAmbiguousScheduleUpdate(schedule []hub.NeuroScheduleEntry, update string) ([]hub.NeuroScheduleEntry, error) {
	ctx := context.Background()

	scheduleJSON, err := json.MarshalIndent(ScheduleJson{
		Schedule: schedule,
		String:   update,
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage("Your job is to take an input JSON object from the user containing a current schedule and an arbitrary string and return the schedule with the change the arbitrary string could be requesting. An example would be if there's a schedule with an event on Thursday at 7pm and the arbitrary string is \"The stream on Thursday will be at 9pm\", then you should modify the schedule to change the time for the event on Thursday. You must respond with the JSON of the schedule array only. Do not wrap it in an object like the user does, only respond with the schedule array you have modified. If changes appear to be targeted at an unknown \"schedule image\" only (e.g. \"The schedule image is incorrect for Ellie's collab\"), then no change is required as the update is focused on correcting a graphic that you are not being given."),
			openai.UserMessage(string(scheduleJSON)),
		},
		Model: openai.ChatModelGPT5_4Nano,
	})
	if err != nil {
		return nil, err
	}
	var parsedSchedule []hub.NeuroScheduleEntry
	err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &parsedSchedule)
	if err != nil {
		return nil, err
	}

	return parsedSchedule, nil
}
