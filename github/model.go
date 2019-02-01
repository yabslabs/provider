package main

type Repos []*Repo

type Repo struct {
	Name string
}

type Repositories struct {
	Repositories []string `json:"repositories"`
}

type Migrations struct {
	ID    int    `json:"id"`
	State string `json:"state"`
}
