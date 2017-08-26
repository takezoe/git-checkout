package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"os/exec"
	"regexp"
	"strings"
)

var s []prompt.Suggest

func completer(in prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	out, _ := exec.Command("git", "branch", "-a").Output()
	lines := strings.Split(string(out), "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Index(line, "*") == 0 {
			s = append(s, prompt.Suggest{Text: line[2:len(line)], Description: "*"})
		} else {
			s = append(s, prompt.Suggest{Text: line})
		}
	}

	in := prompt.Input("branch: ", completer,
		prompt.OptionTitle("git checkout"),
		prompt.OptionPrefixTextColor(prompt.Blue))

	r := regexp.MustCompile(`remotes/(.*)/(.*)`)
	result := r.FindAllStringSubmatch(in, -1)
	if len(result) == 0 {
		out, _ := exec.Command("git", "checkout", in).Output()
		fmt.Println(string(out))
	} else {
		out, _ := exec.Command("git", "checkout", "-b", result[0][2], result[0][1]+"/"+result[0][2]).Output()
		fmt.Println(string(out))
	}
}
