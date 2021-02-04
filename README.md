# Discord

## Example

```go
package main

import (
	"github.com/johynpapin/discord"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const (
	TOKEN = "yourDiscordToken"
)

func main() {
	discord := discord.New(TOKEN)

	if err := discord.Open(); err != nil {
		log.WithError(err).Fatal("error opening the connection to Discord")
	}
	defer discord.Close()

	challenge, err := discord.GetNewChallenge()
	if err != nil {
		log.WithError(err).Fatal("error getting a new challenge")
	}

	log.WithField("challenge", challenge).Info("got a new challenge")

	waitForExit()
}

func waitForExit() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
```
