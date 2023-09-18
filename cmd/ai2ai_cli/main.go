package main

import (
	_ "embed"
	"fmt"
	"log"
	"openaigo/openai"
	"openaigo/setup"

	"github.com/fatih/color"
)

//go:embed agent_1_system_prompt.txt
var systemPromptAgent1 string

//go:embed agent_2_system_prompt.txt
var systemPromptAgent2 string

func main() {
	chatClient, _ := setup.Clients()

	agent1 := chatClient.NewChatWithSystemPrompt(systemPromptAgent1)
	agent2 := chatClient.NewChatWithSystemPrompt(systemPromptAgent2)

	startChat("Normie", agent1, "Physicist", agent2)
}

func startChat(agent1name string, agent1 *openai.Chat, agent2name string, agent2 *openai.Chat) {
	message := "Hello, can you help me?"
	agent1.AppendMessage("system", message)

	for i := 0; i < 4; i++ {
		var err error

		c1 := color.New(color.FgRed)
		c1.Printf("[%s] ", agent1name)
		fmt.Println(message)

		message, err = agent2.Send(message)
		if err != nil {
			log.Fatal(err.Error())
		}

		c2 := color.New(color.FgGreen)
		c2.Printf("[%s] ", agent2name)
		fmt.Println(message)

		message, err = agent1.Send(message)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	c1 := color.New(color.FgRed)
	c1.Printf("[%s] ", agent1name)
	fmt.Println(message)
	fmt.Println()

}
