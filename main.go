package main

import (
	"bufio"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"io/ioutil"
	"log"
	"os"
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Hi, wha is yer fav color: ")
	inputColor, _ := reader.ReadString('\n')

	baseColor, _ := colorful.Hex(inputColor)
	palette := gamut.Triadic(baseColor)

	for i, line := range lines {
		if match, _ := regexp.MatchString("machine_color:", line); match {
			lines[i] = fmt.Sprintf(`machine_color: "%s"`, baseColor.Hex())
		} else if match, _ := regexp.MatchString("dir_color:", line); match {
			color, _ := colorful.MakeColor(palette[0])
			lines[i] = fmt.Sprintf(`dir_color: "%s"`, color.Hex())
		} else if match, _ := regexp.MatchString("git_color:", line); match {
			color, _ := colorful.MakeColor(palette[1])
			lines[i] = fmt.Sprintf(`git_color: "%s"`, color.Hex())
		}
	}

	newConfigFile := strings.Join(lines, "\n")
	err = ioutil.WriteFile(configPath, []byte(newConfigFile), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done!")
}
