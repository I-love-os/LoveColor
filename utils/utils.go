package utils

import (
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"gopkg.in/AlecAivazis/survey.v1"
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

	switch selectedColorOption {
	case "Machine Color":
		color := ""
		colorNamePrompt := &survey.Input{
			Message: "Machine Color:",
		}
		survey.AskOne(colorNamePrompt, &color, nil)

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

		if !confirmation {
			EditSelectedColor()
		}
	}

	return colors
}