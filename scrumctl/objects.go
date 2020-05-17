package scrumctl

type Session struct {
	SessionId   string            `json:"sessionId"`
	Stories     map[string]*Story `json:"stories"`
	Users       map[string]*User  `json:"users"`
	LatestStory string            `json:"latestStory"`
}

type Story struct {
	Name  string         `json:"name"`
	Votes map[string]int `json:"votes"`
}

type User struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}
