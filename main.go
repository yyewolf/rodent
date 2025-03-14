package main

import (
	"github.com/joho/godotenv"
	"github.com/yyewolf/rodent/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
