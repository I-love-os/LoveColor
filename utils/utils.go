package utils

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"strconv"
	"strings"
)

func GetHTMLcolor(colorName string) string {
	if color, ok := colornames.Map[strings.ToLower(colorName)]; ok {
		c, _ := colorful.MakeColor(color)
		return c.Hex()
	} else {
		panic("Color does not exists!")
	}
	return "#FFFFFF"
}

var 	colors = []string {"#FFFFFF", "#FFFFFF", "#FFFFFF"}

func EditSelectedColor() []string {
	selectedColorOption := ""
	selectColorOptionPrompt := &survey.Select{
		Message: "Which color do you want to edit?",
		Options: []string{"Machine Color", "Dir Color", "Git Color", "Finish"},
	}
	survey.AskOne(selectColorOptionPrompt, &selectedColorOption, nil)

	CheckIfEmpty(&selectedColorOption, true, Survey{Prompt:selectColorOptionPrompt, Stype: "select"})

	switch selectedColorOption {
	case "Machine Color":
		color := ""
		colorNamePrompt := &survey.Input{
			Message: "Machine Color:",
		}
		survey.AskOne(colorNamePrompt, &color, nil)

		if CheckIfEmpty(&color, false) {
			fmt.Println(`Color can not be empty. Using default value "#FFFFFF"`)
			color = "#FFFFFF"
		}

		if !strings.HasPrefix(color, "#") {
			 color = GetHTMLcolor(color)
		}

		colors[0] = color
		
		EditSelectedColor()
	case "Dir Color":
		color := ""
		colorNamePrompt := &survey.Input{
			Message: "Dir Color:",
		}
		survey.AskOne(colorNamePrompt, &color, nil)

		if !strings.HasPrefix(color, "#") {
			color = GetHTMLcolor(color)
		}

		colors[1] = color
		
		EditSelectedColor()
	case "Git Color":
		color := ""
		colorNamePrompt := &survey.Input{
			Message: "Git Color:",
		}
		survey.AskOne(colorNamePrompt, &color, nil)

		if !strings.HasPrefix(color, "#") {
			color = GetHTMLcolor(color)
		}

		colors[2] = color
		
		EditSelectedColor()
	case "Finish":
		confirmation := false
		confirmationPrompt := &survey.Confirm{
			Message: "Do you want that color scheme?",
		}
		survey.AskOne(confirmationPrompt, &confirmation, nil)

		tempConfirm := strconv.FormatBool(confirmation)

		CheckIfEmpty(&tempConfirm, true, Survey{Prompt:selectColorOptionPrompt, Stype: "confirm"})

		confirmation, err := strconv.ParseBool(tempConfirm)

		if err != nil {
			panic(err)
		}

		if !confirmation {
			EditSelectedColor()
		}
	}

	return colors
}

type Survey struct {
	Prompt survey.Prompt
	Stype  string
}

func CheckIfEmpty(v *string, exit bool, s ...Survey) bool {
	if *v == "" || *v == " " {
		if exit {
			confirmation := false
			confirmationPrompt := &survey.Confirm{
				Message: "Do you really want to exit? It can corrupt your scheme.conf file. Exit?",
			}
			survey.AskOne(confirmationPrompt, &confirmation, nil)

			if s == nil {
				panic("Survey should not be empty")
			}

			if confirmation {
				os.Exit(1)
			} else {
				switch s[0].Stype {
				case "confirm":
					o := false
					survey.AskOne(s[0].Prompt, &o, nil)
					*v = strconv.FormatBool(o)
					return true
				case "select":
					survey.AskOne(s[0].Prompt, v, nil)
				case "input":
					survey.AskOne(s[0].Prompt, v, nil)
				default:
					panic(fmt.Sprintf(`sType "%s" does not exist`, s[0].Stype))
				}

				CheckIfEmpty(v, true, Survey{Prompt:s[0].Prompt, Stype:s[0].Stype})
			}
		}
		return true
	}
	return false
}