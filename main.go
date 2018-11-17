package main

import (
	"LoveColor/utils"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"gopkg.in/AlecAivazis/survey.v1"
	"image/color"
	"io/ioutil"
	"log"
	"os/user"
	"regexp"
	"strings"
)

func main() {
	currentUser, err := user.Current()

	if err != nil {
		panic(err)
	}

	configPath := fmt.Sprintf(`/home/%s/.config/LoveShell/LoveShell.conf`, currentUser.Username)
	configFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(configFile), "\n")

	initType := ""
	initTypePrompt := &survey.Select{
		Message: "Select how would you like to create your color scheme:",
		Options: []string{"Automatically", "Set Manually"},
	}
	survey.AskOne(initTypePrompt, &initType, nil)

	var palette []color.Color
	var colors []string

	if initType == "Automatically" {
		inputColor := ""
		askColorPrompt := &survey.Input{Message: "Hi, wha is yer fav color?"}
		survey.AskOne(askColorPrompt, &inputColor, nil)

		if !strings.HasPrefix(inputColor, "#") {
			inputColor = utils.GetHTMLcolor(inputColor)
		}

		baseColor, _ := colorful.Hex(inputColor)

		genType := ""
		genTypePrompt := &survey.Select{
			Message: "How do you want to generate your color scheme?",
			Options: []string{"Triadic", "Monochromatic", "Shades"},
		}
		survey.AskOne(genTypePrompt, &genType, nil)

		switch genType {
		case "Triadic":
			palette = gamut.Triadic(baseColor)
		case "Monochromatic":
			palette = gamut.Monochromatic(baseColor, 2)
		case "Shades":
			palette = gamut.Shades(baseColor, 2)
		default:
			panic(fmt.Sprintf(`Option "%s" does not exist`, genType))
		}

		colors = append(colors, baseColor.Hex())

		for _, v := range palette {
			c, _ := colorful.MakeColor(v)
			colors = append(colors, c.Hex())
		}
	} else {
		colors = utils.EditSelectedColor()
	}


	for i, line := range lines {
		if match, _ := regexp.MatchString("machine_color:", line); match {
			lines[i] = fmt.Sprintf(`machine_color: "%s"`, colors[0])
		} else if match, _ := regexp.MatchString("dir_color:", line); match {
			lines[i] = fmt.Sprintf(`dir_color: "%s"`, colors[1])
		} else if match, _ := regexp.MatchString("git_color:", line); match {
			lines[i] = fmt.Sprintf(`git_color: "%s"`, colors[2])
		}
	}

	newConfigFile := strings.Join(lines, "\n")
	err = ioutil.WriteFile(configPath, []byte(newConfigFile), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done!")
}
