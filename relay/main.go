package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

type PlantStatus struct {
	Plant  string
	Status string
}

// Relay plant status to Telegram.
// Listens on HTTP on purpose.
// Example request:
// GET /notify?plant=front-porch-mint&status=thirsty&secret=...
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	secret := os.Getenv("PLANT_SECRET")

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("secret") != secret {
			log.Printf("Failed: unauthorized with secret %s", r.URL.Query().Get("secret"))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		whichPlant := r.URL.Query().Get("plant")
		if whichPlant == "" {
			whichPlant = "Unknown plant"
		}

		plantStatus := r.URL.Query().Get("status")
		if plantStatus == "" {
			plantStatus = "unknown status"
		}

		status := PlantStatus{
			Plant:  whichPlant,
			Status: plantStatus,
		}

		endpoint := fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage",
			token,
		)

		tmpl := template.Must(template.New("").Parse(`🪴 {{ .Plant }} is {{ .Status }}`))
		writer := bytes.NewBufferString("")
		_ = tmpl.Execute(writer, status)

		text := writer.String()
		resp, err := http.PostForm(endpoint, url.Values{
			"chat_id": {chatID},
			"text":    {text},
		})
		if err != nil {
			log.Printf("Failed: telegram request failed (%v)", err)
			http.Error(w, "telegram request failed", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			log.Printf("Failed: telegram rejected request (%d)", resp.StatusCode)
			http.Error(w, "telegram rejected request", http.StatusBadGateway)
			return
		}
		log.Printf("Success: sent %s to telegram", text)

		fmt.Fprintln(w, "ok")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🪴 Listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
