package server

type Repo struct {
	Url string `json:"url"`
	Directory string `json:"directory"`
	Hidden bool `json:"hidden"`
}

type GitsbyConfig struct {
	Landing bool `json:"landing"`
	Repos []Repo `json:"repos"`
}
