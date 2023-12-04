package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MohammadBnei/go-openai-cli/service"
	"github.com/MohammadBnei/go-openai-cli/ui"
	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
)

func AddFileCommand(commandMap map[string]func(*PromptConfig) error) {
	commandMap["save"] = func(pc *PromptConfig) error {
		assistantRole := service.RoleAssistant
		lastMessage := pc.ChatMessages.LastMessage(&assistantRole)
		if lastMessage == nil {
			return errors.New("no assistant message found")
		}
		return ui.SaveToFile([]byte(lastMessage.Content), "")
	}

	commandMap["file"] = func(pc *PromptConfig) error {
		fileContents, err := ui.FileSelectionFzf()
		if err != nil {
			return err
		}
		for _, fileContent := range fileContents {
			pc.ChatMessages.AddMessage(fileContent, service.RoleUser)
		}
		return nil
	}
}

func AddConfigCommand(commandMap map[string]func(*PromptConfig) error) {
	commandMap["markdown"] = func(pc *PromptConfig) error {
		pc.MdMode = !pc.MdMode
		return nil
	}
}

func AddSystemCommand(commandMap map[string]func(*PromptConfig) error) {
	commandMap["list"] = func(pc *PromptConfig) error {
		messages, err := ui.SelectSystemCommand()
		if err != nil {
			return err
		}
		for _, message := range messages {
			pc.ChatMessages.AddMessage(message, service.RoleAssistant)
		}
		return nil
	}

	commandMap["d-list"] = func(pc *PromptConfig) error {
		return ui.DeleteSystemCommand()
	}

	commandMap["system"] = func(pc *PromptConfig) error {
		message, err := ui.SendAsSystem()
		if err != nil {
			return err
		}
		pc.ChatMessages.AddMessage(message, service.RoleAssistant)
		return nil
	}

	commandMap["filter"] = func(pc *PromptConfig) error {
		messageIds, err := ui.FilterMessages(pc.ChatMessages.Messages)
		if err != nil {
			return err
		}

		for _, id := range messageIds {
			_err := pc.ChatMessages.DeleteMessage(id)
			if _err != nil {
				err = errors.Join(err, _err)
			}
		}

		return err
	}

	commandMap["cli-clear"] = func(pc *PromptConfig) error {
		ui.ClearTerminal()
		return nil
	}

	commandMap["reuse"] = func(pc *PromptConfig) error {
		message, err := ui.ReuseMessage(pc.ChatMessages.Messages)
		if err != nil {
			return err
		}
		pc.PreviousPrompt = message
		return nil
	}

	commandMap["responses"] = func(pc *PromptConfig) error {
		_, err := ui.ShowPreviousMessage(pc.ChatMessages.Messages, pc.MdMode)
		return err
	}

	commandMap["default"] = func(pc *PromptConfig) error {
		commandToAdd, err := ui.SetSystemDefault(false)
		if err != nil {
			return err
		}
		for _, command := range commandToAdd {
			pc.ChatMessages.AddMessage(command, service.RoleAssistant)
		}
		return nil
	}
	commandMap["d-default"] = func(pc *PromptConfig) error {
		commandToAdd, err := ui.SetSystemDefault(true)
		if err != nil {
			return err
		}
		for _, command := range commandToAdd {
			pc.ChatMessages.AddMessage(command, service.RoleAssistant)
		}
		return nil
	}

	commandMap["copy"] = func(pc *PromptConfig) error {
		assistantMessages, _ := pc.ChatMessages.FilterMessages(service.RoleAssistant)
		if len(assistantMessages) < 1 {
			return errors.New("no messages to copy")
		}

		clipboard.WriteAll(assistantMessages[len(assistantMessages)-1].Content)
		fmt.Println("copied to clipboard")
		return nil
	}

	commandMap["clear"] = func(pc *PromptConfig) error {
		pc.ChatMessages.ClearMessages()
		fmt.Println("cleared messages")
		return nil
	}
}

func AddImageCommand(commandMap map[string]func(*PromptConfig) error) {
	commandMap["image"] = func(cfg *PromptConfig) error {
		return ui.AskForImage()
	}
}

func AddHuggingFaceCommand(commandMap map[string]func(*PromptConfig) error) {
	commandMap["mask"] = func(cfg *PromptConfig) error {
		maskPrompt := promptui.Prompt{
			Label: "Write a sentance with the character !! as the token to replace",
		}
		pr, err := maskPrompt.Run()
		if err != nil {
			return err
		}
		result, err := service.Mask(strings.Replace(pr, "!!", "[MASK]", -1))
		if err != nil {
			return err
		}
		fmt.Println("Result : ", result)

		return nil
	}
}
