package main

import (
	"bufio"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"os"
	"os/user"
)


func main() {
	currentUser, err := user.Current()

	if err != nil {
		panic(err)
	}

	configPath := fmt.Sprintf(`/home/%s/.config/LoveShell/LoveShell.conf`, currentUser.Username)
	configFile, err := os.Open(configPath)

	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(configFile)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	color, _ := colorful.Hex("#8A2BE2")
	pallette := gamut.Monochromatic(color, 3)
	//fmt.Println(pallette)
	for _, v := range pallette {
		c, _ := colorful.MakeColor(v)
		fmt.Println(c.Hex())
	}
}
