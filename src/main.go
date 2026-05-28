package main

import (
	"fmt"
	"neuro-schedule-api/src/discord"
	"neuro-schedule-api/src/http"
	"neuro-schedule-api/src/hub"
	"neuro-schedule-api/src/parser"
	"os"
	"os/signal"
	"syscall"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	h := hub.New()
	httpServer := http.New(h)
	parser.ConnectToOpenAI(os.Getenv("OPENAI_KEY"))
	err := discord.StartDiscordListener(os.Getenv("DISCORD_BOT_TOKEN"), os.Getenv("DISCORD_CHANNEL_ID"), h)
	if err != nil {
		fmt.Println("Encountered an error while trying to start the Discord listener:", err)
		panic(err)
	}

	err = httpServer.Start()
	if err != nil {
		fmt.Println("Encountered an error while trying to start the HTTP server:", err)
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = discord.StopDiscordListener()
	if err != nil {
		fmt.Println("Encountered an error while trying to stop the Discord listener:", err)
		return
	}
}
