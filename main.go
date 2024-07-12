package main

import "github.com/9-Realms-Dev/go_nvm/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
