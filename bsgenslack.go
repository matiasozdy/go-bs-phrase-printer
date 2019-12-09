package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "github.com/nlopes/slack"
  "regexp"
  "strings"
  "os"
)

type Phrase struct {
  Phrase string
}

func getenv(name string) string {
  v := os.Getenv(name)
  if v == "" {
    panic("missing required environment variable " + name)
  }
  return v
}

func main() {

  //Slack init
  token := getenv("SLACKTOKEN")
  api := slack.New(token)
  rtm := api.NewRTM()
  go rtm.ManageConnection()

  //Slack loop
Loop:
  for {
    select {
    case msg := <-rtm.IncomingEvents:
      fmt.Print("\nEvent Received:")
      switch ev := msg.Data.(type) {

      case *slack.MessageEvent:
        info := rtm.GetInfo()

        text := ev.Text
        text = strings.TrimSpace(text)
        text = strings.ToLower(text)
        botName := strings.ToLower(info.User.ID)

        matched, _ := regexp.MatchString("<@"+botName+"> inspire me", text)
        fmt.Print("text: "+ text + "user:" + botName)
        if ev.User != info.User.ID && matched {
          //Get random buzzwords
          var client http.Client
          var url = "https://corporatebs-generator.sameerkumar.website"
          resp, err := client.Get(url)
          if err != nil {
            log.Fatal(err)
          }
          defer resp.Body.Close()

          if resp.StatusCode == http.StatusOK {
            bodyBytes, err := ioutil.ReadAll(resp.Body)
            if err != nil {
              log.Fatal(err)
            }
            var phrase1 Phrase
            jsonErr := json.Unmarshal(bodyBytes, &phrase1)
            if jsonErr != nil {
              log.Fatal(jsonErr)
            }
            rtm.SendMessage(rtm.NewOutgoingMessage(phrase1.Phrase, ev.Channel))
          }
        }

      case *slack.RTMError:
        fmt.Printf("Error: %s\n", ev.Error())

      case *slack.InvalidAuthEvent:
        fmt.Printf("Invalid credentials")
        break Loop

      default:
        // Take no action
      }
    }
  }
// end slack loop
}
