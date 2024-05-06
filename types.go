package main

import "time"

type RepoData struct {
	Username    string
	Repo        string
	GithubToken string
}

type GitHubVariable interface {
	Edit(int, string) error
}

type RepoVar struct {
	Id        int
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VarResponse struct {
	Variables  []*RepoVar `json:"variables"`
	TotalCount int        `json:"total_count"`
}

type RepoSecret struct {
	Id        int
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SecretResponse struct {
	Secrets    []*RepoSecret `json:"secrets"`
	TotalCount int           `json:"total_count"`
}

type VarDB struct {
	repo_data     *RepoData
	vars          []*RepoVar
	secrets       []*RepoSecret
	currentSecret string
}

type EncryptedRequestBody struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
}
