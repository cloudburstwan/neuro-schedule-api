package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type Hub struct {
	mu sync.RWMutex

	schedule []NeuroScheduleEntry
}

type NeuroScheduleEntry struct {
	Timestamp string   `json:"timestamp"`
	Title     string   `json:"title,omitempty"`
	Live      bool     `json:"live"`
	Streamers []string `json:"streamers,omitempty"`
}

func New() *Hub {
	file, err := os.OpenFile("./data/schedule.json", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Encountered an error while opening schedule.json: ", err)
		panic(err)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Encountered an error while reading schedule.json: ", err)
		panic(err)
	}

	if len(bytes) == 0 {
		// Nothing in it, just initiate empty schedule
		err = file.Close()
		if err != nil {
			fmt.Println("Encountered an error while closing schedule.json: ", err)
			// Not really the biggest deal on the planet
		}

		return &Hub{
			schedule: []NeuroScheduleEntry{},
		}
	}

	var schedule []NeuroScheduleEntry
	err = json.Unmarshal(bytes, &schedule)
	if err != nil {
		fmt.Println("Encountered an error while parsing schedule.json: ", err)
		panic(err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println("Encountered an error while closing schedule.json: ", err)
		// Not really the biggest deal on the planet
	}

	return &Hub{
		schedule: schedule,
	}
}

func (h *Hub) Broadcast(newSchedule []NeuroScheduleEntry) error {
	h.mu.Lock()
	h.schedule = newSchedule
	h.mu.Unlock()

	// TODO: go channel support for real-time

	scheduleJSON, err := json.Marshal(h.schedule)
	if err != nil {
		return err
	}
	err = os.WriteFile("./data/schedule.json", scheduleJSON, 0644)
	if err != nil {
		fmt.Println("Encountered an error while writing schedule.json: ", err)
		panic(err)
	}
	return nil
}

func (h *Hub) GetSchedule() []NeuroScheduleEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.schedule
}
