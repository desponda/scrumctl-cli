package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"net/http"
	"net/url"
	"scrumctl-cli/scrumctl"
	"strconv"
	"time"
)

func main() {
	prompt := promptui.Select{
		Label: "Choose an option",
		Items: []string{"Create a session", "Join a session"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Something went wrong")
		return
	}
	u, _ := url.Parse("http://localhost:3000")
	c := &scrumctl.Client{
		BaseURL:    u,
		UserAgent:  "",
		HttpClient: http.DefaultClient,
	}

	promptName := promptui.Prompt{
		Label: "Name",
	}
	creator := false
	var s scrumctl.Session
	un, _ := promptName.Run()
	if result == "Create a session" {
		fmt.Printf("Creating a session...\n")
		s, _ = c.CreateSession(un)
		fmt.Printf("Created session %v\n", s.SessionId)
		creator = true
	} else {
		promptName = promptui.Prompt{
			Label: "Session ID",
		}
		sessionId, _ := promptName.Run()
		s = c.JoinSession(sessionId, un)
		fmt.Printf("Joined session %v\n", s.SessionId)
	}
	for {
		var story scrumctl.Story
		if creator {
			promptName = promptui.Prompt{
				Label: "Story ID",
			}
			sn, _ := promptName.Run()
			story, _ = c.CreateStory(sn, s.SessionId)
		} else {
			story = *s.Stories[s.LatestStory]
		}
		validate := func(input string) error {
			_, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return errors.New("invalid number")
			}
			return nil
		}
		storyPrompt := fmt.Sprintf("Vote on story %q", story.Name)
		promptName = promptui.Prompt{
			Label:    storyPrompt,
			Validate: validate,
		}
		r, _ := promptName.Run()
		vote, _ := strconv.ParseInt(r, 10, 64)
		_ = c.CastVote(s.SessionId, story.Name, un, int(vote))
		s, _ = c.FindSession(s.SessionId)
		for !allUsersVoted(s) {
			time.Sleep(3 * time.Second)
			s, _ = c.FindSession(s.SessionId)
			fmt.Printf("Waiting for all users to vote...\n")
		}
		fmt.Printf("All users have voted! Results:\n")
		for u, v := range s.Stories[s.LatestStory].Votes {
			fmt.Printf("%v: %v\n", u, v)
		}
	}
}

func allUsersVoted(s scrumctl.Session) bool {
	for user, _ := range s.Users {
		if _, ok := s.Stories[s.LatestStory].Votes[user]; !ok {
			return false
		}
	}
	return true
}
