package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
)

func main() {
	ttsClient := setup.Clients().TTS()
	ttsClient.SetHD(true)

	input := strings.Join(os.Args[1:], " ")
	if input == "" {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		input = string(stdin)
		if input == "" {
			log.Fatal("Error: no input")
		}
	}

	audio, err := ttsClient.GenerateAudio(input)
	if err != nil {
		log.Fatal("Error: couldn't generate audio")
	}

	play(audio)
}

func play(audio *openai.Audio) {
	if !commandExists("ffplay") {
		log.Fatal("Error: need ffplay to play audio")
	}

	cmd := exec.Command("ffplay", "-v", "0", "-nodisp", "-autoexit", "-")
	cmd.Stdin = bytes.NewReader(audio.Data)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
