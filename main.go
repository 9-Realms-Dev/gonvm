package main

import "github.com/9-Realms-Dev/gonvm/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
