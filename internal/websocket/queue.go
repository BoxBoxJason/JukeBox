package websocket

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
)

const (
	INITIAL_PROMPT            = "Please play a music that factors in all of the parameters in between the separator"
	SEPARATOR_AVAILABLE_CHARS = "&~\"'-_+=*^%$#@!`"
	MINIMUM_SEPARATOR_LENGTH  = 5
	MAXIMUM_SEPARATOR_LENGTH  = 15
	RETRY_DELAY               = 5 * time.Second
	PROMPT_INTERVAL           = 30 * time.Second
)

var (
	current_message_stack = []string{}
	mu_message_stack      = sync.RWMutex{}
	MUSIC_GENERATOR_URL   = os.Getenv("MUSIC_GENERATOR_URL")
	prompt_lock           = sync.Mutex{}
)

type MusicGeneratorRequest struct {
	Request  string `json:"request"`
	Duration int    `json:"duration,omitempty"`
}

func init() {
	go watchPrompts()
}

// addMessage adds a message to the current message stack
func addMessage(message string) {
	mu_message_stack.Lock()
	current_message_stack = append(current_message_stack, message)
	mu_message_stack.Unlock()
	if len(current_message_stack) >= 10 {
		sendPrompt()
	}
}

// createRandomSeparator creates a random separator
func createRandomSeparator() string {
	length := rand.Intn(MAXIMUM_SEPARATOR_LENGTH-MINIMUM_SEPARATOR_LENGTH+1) + MINIMUM_SEPARATOR_LENGTH
	separator := make([]byte, length)
	for i := range separator {
		separator[i] = SEPARATOR_AVAILABLE_CHARS[rand.Intn(len(SEPARATOR_AVAILABLE_CHARS))]
	}
	return string(separator)
}

// sendPrompt sends a retryable prompt to the music generator
func sendPrompt() error {
	random_separator := createRandomSeparator()
	mu_message_stack.RLock()
	raw_request := INITIAL_PROMPT + "\n" + random_separator + strings.Join(current_message_stack, "\n") + random_separator
	mu_message_stack.RUnlock()
	request := MusicGeneratorRequest{
		Request:  raw_request,
		Duration: 10,
	}

	prompt_lock.Lock()
	defer prompt_lock.Unlock()

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	logger.Debug("Sending prompt to music generator")
	err = httputils.RetryRequest(MUSIC_GENERATOR_URL, jsonData, 5, 3)
	if err != nil {
		return err
	}

	emptyMessageStack()
	return nil
}

// watchPrompts watches the current message stack and sends prompts to the music generator
// when either the stack reaches 10 messages OR the interval is reached
func watchPrompts() {
	ticker := time.NewTicker(PROMPT_INTERVAL)
	defer ticker.Stop()
	for range ticker.C {
		mu_message_stack.RLock()
		if len(current_message_stack) > 0 {
			mu_message_stack.RUnlock()
			err := sendPrompt()
			if err != nil {
				logger.Error("Failed to send prompt to music generator", err)
			}
		} else {
			mu_message_stack.RUnlock()
		}
	}
}

// emptyMessageStack empties the message stack
func emptyMessageStack() {
	mu_message_stack.Lock()
	defer mu_message_stack.Unlock()
	current_message_stack = []string{}
}
