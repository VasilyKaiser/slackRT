package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
)

var (
	OAUTH_TOKEN    string
	api            *slack.Client
	CHANNEL_ID     string
	SIGNING_SECRET string
	PORT           string
)

func init() {
	parseYamlFile()
	api = slack.New(OAUTH_TOKEN)
}

type configs struct {
	OAUTH_TOKEN    string `yaml:"OAUTH_TOKEN"`
	CHANNEL_ID     string `yaml:"CHANNEL_ID"`
	SIGNING_SECRET string `yaml:"SIGNING_SECRET"`
	PORT           string `yaml:"PORT"`
}

func (c *configs) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func parseYamlFile() {
	data, err := os.ReadFile("config.yaml")
	checkErr(err)
	var config configs
	err = config.Parse(data)
	checkErr(err)
	OAUTH_TOKEN = config.OAUTH_TOKEN
	CHANNEL_ID = config.CHANNEL_ID
	SIGNING_SECRET = config.SIGNING_SECRET
	PORT = config.PORT

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func deleteAllMessages() {
	info, _ := api.GetConversationHistory(&slack.GetConversationHistoryParameters{ChannelID: CHANNEL_ID})
	for _, message := range info.Messages {
		log.Println(api.DeleteMessage(CHANNEL_ID, message.Msg.Timestamp))
	}
}

func sendMessage(message string) {
	_, _, err := api.PostMessage(
		CHANNEL_ID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(false),
	)
	checkErr(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, SIGNING_SECRET)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println(s.Command)
	switch s.Command {
	case "/sh":
		if s.Text == "del" {
			deleteAllMessages()

		} else if s.Text != "" {
			command := strings.Split(s.Text, " ")
			var out []byte
			var err error
			if len(command) > 1 {
				out, err = exec.Command(command[0], command[1:]...).CombinedOutput()
			} else {
				out, err = exec.Command(command[0]).CombinedOutput()
			}
			if len(out) > 0 {
				message := string(out)
				if err != nil {
					message += "\n" + err.Error()
				}
				sendMessage(message)
			} else {
				sendMessage("Executed, but didn't get response from command yet.")
			}
		} else {
			sendMessage("I need an argument to the `/sh` command!\nExample: `/sh del`")
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("[INFO] Server listening on port: " + PORT)

	http.ListenAndServe(":"+PORT, nil)
}
