package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "strconv"
	"github.com/jackc/pgx/v4/pgxpool"
	tb "gopkg.in/telebot.v3"
)

var bot *tb.Bot
var postgres *pgxpool.Pool

// New initializes the bot and database pool for the telegram package
func New(b *tb.Bot, p *pgxpool.Pool) {
	bot = b
	postgres = p
}

// HandleWebhook handles incoming webhook requests from Telegram
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var update tb.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if update.Message == nil {
		http.Error(w, "No message found", http.StatusBadRequest)
		return
	}

	message := update.Message
	username := message.Sender.Username
	if username == "" {
		username = message.Sender.FirstName
	}
	name := message.Sender.FirstName
	chatID := int(message.Chat.ID) // Convert chatID to int
	text := message.Text

	log.Printf("Chat ID: %d, Text: %s", chatID, text)

	if text == "/start" {
		log.Println("Received /start command")
    
    exists, err := isUserExists(chatID)

    if err != nil {
      log.Printf("Failed to check if user exists: %v", err)
      return
    }

    if exists {
      log.Println("User already exists")
      w.WriteHeader(http.StatusOK)
      return
    }

		err = createUser(username, chatID, name, "")
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			return
	  }	
    button := tb.InlineButton{
			Text: "Открыть веб-приложение",
			WebApp: &tb.WebApp{
				URL: fmt.Sprintf("https://prichal.weble.kz/home/%s", strconv.Itoa(chatID)),
			},
		}
		// Create a keyboard with a web app button
		keyboard := &tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{
				{button},
			},
		}

		bot.Send(message.Chat, "Click to open the web application", &tb.SendOptions{
			ReplyMarkup: keyboard,
		})

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// createUser inserts a new user into the database
func createUser(username string, chatID int, name string, imageURL string) error {
	const query = `
		INSERT INTO users (username, chat_id, name, image_url)
		VALUES ($1, $2, $3, $4)
	`
	_, err := postgres.Exec(context.Background(), query, username, chatID, name, imageURL)
	return err
}

func isUserExists(chatID int) (bool, error) {
  const query = `
    SELECT EXISTS(SELECT 1 FROM users WHERE chat_id = $1)
  `
  var exists bool
  err := postgres.QueryRow(context.Background(), query, chatID).Scan(&exists)
  if err != nil {
    return false, err
  }
  return exists, nil
}
