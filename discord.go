package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const (
	CHANNEL_ID       = "806882578086887445"
	SAN_JOSE_ID      = "805794737252597792"
	CHALLENGE_PREFIX = "challenge "
)

type Discord struct {
	token   string
	session *discordgo.Session

	lastRequestMessageId string
	challenge            chan string
}

func New(token string) *Discord {
	return &Discord{token: token, challenge: make(chan string)}
}

func (discord *Discord) Open() error {
	var err error
	discord.session, err = discordgo.New("Bot " + discord.token)
	if err != nil {
		return fmt.Errorf("creating the Discord session: %w", err)
	}

	discord.session.AddHandler(discord.handleMessageCreate)

	if err = discord.session.Open(); err != nil {
		return fmt.Errorf("opening the connection to Discord: %w", err)
	}

	return nil
}

func (discord *Discord) Close() error {
	if err := discord.session.Close(); err != nil {
		return fmt.Errorf("closing the connection to Discord: %w", err)
	}

	close(discord.challenge)

	return nil
}

func (discord *Discord) GetNewChallenge() (string, error) {
	requestMessage, err := discord.session.ChannelMessageSend(CHANNEL_ID, "I want a new challenge.")
	if err != nil {
		return "", fmt.Errorf("sending the challenge request: %w", err)
	}

	discord.lastRequestMessageId = requestMessage.ID

	challenge, ok := <-discord.challenge
	if !ok {
		return "", fmt.Errorf("closed")
	}

	return challenge, nil
}

func (discord *Discord) SendGuess(guess string) error {
	if _, err := discord.session.ChannelMessageSend(CHANNEL_ID, guess); err != nil {
		return fmt.Errorf("sending the guess: %w", err)
	}

	return nil
}

func (discord *Discord) handleMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	author := message.Message.Author
	content := message.Message.Content

	// If the message doesn't come from San José, we ignore it.
	if author.ID != SAN_JOSE_ID {
		return
	}

	if message.Message.ChannelID != CHANNEL_ID {
		return
	}

	// If the message doesn't starts with the challenge prefix, we ignore it.
	if !strings.HasPrefix(content, CHALLENGE_PREFIX) {
		return
	}

	// If the message isn't a reply, we ignore it.
	parentMessageReference := message.MessageReference
	if parentMessageReference == nil {
		return
	}

	if parentMessageReference.MessageID != discord.lastRequestMessageId {
		return
	}

	challenge := strings.TrimPrefix(content, CHALLENGE_PREFIX)

	discord.challenge <- challenge
}
