package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
	"github.com/fatih/color"
)

//go:embed agent_1_system_prompt.txt
var systemPromptAgent1 string

//go:embed agent_2_system_prompt.txt
var systemPromptAgent2 string

func main() {
	chatClient := setup.Clients().Chat()

	agent1 := chatClient.NewChatWithSystemPrompt(systemPromptAgent1)
	agent2 := chatClient.NewChatWithSystemPrompt(systemPromptAgent2)

	startChat("Normie", agent1, "Physicist", agent2)
}

func startChat(agent1name string, agent1 *openai.Chat, agent2name string, agent2 *openai.Chat) {
	redFmt := color.New(color.FgRed)
	greenFmt := color.New(color.FgGreen)

	message := "Hello, can you help me?"
	agent1.AppendMessage("system", message)

	redFmt.Printf("[%s] ", agent1name)
	fmt.Println(message)

	message, err := agent2.Send(message)
	if err != nil {
		log.Fatal(err.Error())
	}

	greenFmt.Printf("[%s] ", agent2name)
	fmt.Println(message)

	for i := 0; i < 4; i++ {

		message, err = agent1.Send(message)
		if err != nil {
			log.Fatal(err.Error())
		}

		redFmt.Printf("[%s] ", agent1name)
		fmt.Println(message)

		message, err = agent2.Send(message)
		if err != nil {
			log.Fatal(err.Error())
		}

		greenFmt.Printf("[%s] ", agent2name)
		fmt.Println(message)
	}

}
