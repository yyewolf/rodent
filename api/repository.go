package api

import "github.com/go-fuego/fuego"

type Repository interface {
	Group() string
	Register(server *fuego.Server)
}
