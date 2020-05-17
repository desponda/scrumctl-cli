package main

import (
	"errors"
	"fmt"
	"github.com/desponda/scrumctl-cli/scrumctl"
	"github.com/manifoldco/promptui"
	"net/url"
	"strconv"
	"time"
)

func main() {
	prompt := promptui.Prompt{
		Label: "Username",
	}
	un, _ := prompt.Run()

	host, _ := url.Parse("https://scrumctl.dev")
	c := scrumctl.NewClient(host)
	creator := false
	var s scrumctl.Session
	s, creator = initializeSession(c, un)

	for {

		story := getActiveStory(c, creator, s)
		vote, _ := getVote(&story)
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

func getVote(story *scrumctl.Story) (int64, error) {
	validate := func(input string) error {
		_, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}
	storyPrompt := fmt.Sprintf("Vote on story %q", story.Name)
	prompt := promptui.Prompt{
		Label:    storyPrompt,
		Validate: validate,
	}
	r, _ := prompt.Run()
	return strconv.ParseInt(r, 10, 64)
}

func getActiveStory(c *scrumctl.Client, creator bool, s scrumctl.Session) (story scrumctl.Story) {
	if creator {
		prompt := promptui.Prompt{
			Label: "Story ID",
		}
		sn, _ := prompt.Run()
		story, _ = c.CreateStory(sn, s.SessionId)
	} else {
		story = *s.Stories[s.LatestStory]
	}
	return story
}

func initializeSession(c *scrumctl.Client, un string) (s scrumctl.Session, creator bool) {
	promptSelect := promptui.Select{
		Label: "Choose an option",
		Items: []string{"Create a session", "Join a session"},
	}
	_, result, err := promptSelect.Run()
	if err != nil {
		fmt.Printf("Something went wrong")
		return
	}
	if result == "Create a session" {
		fmt.Printf("Creating a session...\n")
		s, _ = c.CreateSession(un)
		fmt.Printf("Created session %v\n", s.SessionId)
		creator = true
	} else {
		prompt := promptui.Prompt{
			Label: "Session ID",
		}
		sessionId, _ := prompt.Run()
		s = c.JoinSession(sessionId, un)
		fmt.Printf("Joined session %v\n", s.SessionId)
	}
	return s, creator
}

func allUsersVoted(s scrumctl.Session) bool {
	for user, _ := range s.Users {
		if _, ok := s.Stories[s.LatestStory].Votes[user]; !ok {
			return false
		}
	}
	return true
}
