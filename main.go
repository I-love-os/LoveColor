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
	//TODO: Handle ctrl + c

	configPath := fmt.Sprintf(`/home/%s/.config/Love/schemes.conf`, currentUser.Username)
	configFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(configFile), "\n")

	var schemes []string

	for _, line := range lines {
		if match, _ := regexp.MatchString(`\w+:.*?{`, line); match && !strings.HasPrefix(line, "#") {
			replacer := strings.NewReplacer(":", ""," ", "", "{", "", "}", "")
			schemes = append(schemes, replacer.Replace(line))
		}
	}

	schemes = append(schemes, "Add new scheme")

	selectedScheme := ""
	selectedSchemePrompt := &survey.Select{
		Message: "Which scheme d'ye wanna edit?",
		Options: schemes,
	}
	survey.AskOne(selectedSchemePrompt, &selectedScheme, nil)

	if selectedScheme == "Add new scheme" {
		schemeName := ""
		schemeNamePrompt := &survey.Input{Message: "So, how would you like to name this AWESOME scheme?"}
		survey.AskOne(schemeNamePrompt, &schemeName, nil)

		if schemeName == "" || schemeName == " " {
			schemeName := ""
			schemeNamePrompt := &survey.Input{Message: "So AGAIN, how would you like to name this AWESOME scheme (remember you can not type just empty space)?"}
			survey.AskOne(schemeNamePrompt, &schemeName, nil)
		}

		lines = append(lines, " ", fmt.Sprintf(`    %s: {`, schemeName), `      machine_color: "#FFFFFF"`, `      dir_color: "#FFFFFF"`, `      git_color: "#FFFFFF"`, `      git_diff_color: "#f6f4f5"`, `      font_color: "#495049"`, "    }")

		selectedScheme = schemeName
	}

	selectedMode := ""
	selectModePrompt := &survey.Select{
		Message: "Select how would you like to create your color scheme:",
		Options: []string{"Automatically", "Set Manually"},
	}
	survey.AskOne(selectModePrompt, &selectedMode, nil)

	var palette []color.Color
	var colors []string

	if selectedMode == "Automatically" {
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

	isThisASSchemeIwant := false
	for i, line := range lines {
		if match, _ := regexp.MatchString(selectedScheme, line); match {
			isThisASSchemeIwant = true
		}

		if match, _ := regexp.MatchString("}", line); match {
			isThisASSchemeIwant = false
		}

		if !isThisASSchemeIwant {
			continue
		}

		if match, _ := regexp.MatchString("machine_color:", line); match {
			lines[i] = fmt.Sprintf(`      machine_color: "%s"`, colors[0])
		} else if match, _ := regexp.MatchString("dir_color:", line); match {
			lines[i] = fmt.Sprintf(`      dir_color: "%s"`, colors[1])
		} else if match, _ := regexp.MatchString("git_color:", line); match {
			lines[i] = fmt.Sprintf(`      git_color: "%s"`, colors[2])
		}

	}

	newConfigFile := strings.Join(lines, "\n")
	err = ioutil.WriteFile(configPath, []byte(newConfigFile), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done!")
}
