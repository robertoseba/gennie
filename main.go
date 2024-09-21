package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/output/menu"
)

func main() {
	menu := menu.NewMenu("Select a model")
	menu.AddItem("OpenAI", "openai")
	menu.AddItem("ClaudeAI", "claudeai")
	res := menu.Display()
	fmt.Println(res)
	os.Exit(0)
	// client := httpclient.NewClient()

	// inputOptions := app.ParseCliOptions()

	// model := models.NewModel(inputOptions.Model, client)

	// res, err := model.Ask(inputOptions.Question, nil)

	// if err != nil {
	// 	fmt.Fprint(os.Stderr, err)
	// 	fmt.Println()
	// 	os.Exit(1)
	// 	return
	// }

	// fmt.Printf(res.Answer.Content)

	// os.Exit(0)
}
