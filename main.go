package main

import (
	"github.com/joho/godotenv"
	"github.com/yyewolf/rodent/cmd"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
