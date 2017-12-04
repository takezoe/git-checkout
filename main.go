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
	return prompt.FilterContains(s, in.GetWordBeforeCursor(), true)
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
		checkout(in, func() {
			out, _ := exec.Command("git", "checkout", in).Output()
			fmt.Println(string(out))
		})
	} else {
		checkout(result[0][2], func() {
			out, _ := exec.Command("git", "checkout", "-b", result[0][2], result[0][1]+"/"+result[0][2]).Output()
			fmt.Println(string(out))
		})
	}
}

func checkout(branch string, fn func()) {
	if branch != currentBranch() {
		isDirty := isWorkingTreeDirty()
		if !isDirty {
			exec.Command("git", "stash", "save")
		}
		fn()
		if !isDirty {
			exec.Command("git", "stash", "pop")
		}
	}
}

func currentBranch() string {
	out, _ := exec.Command("git", "branch").Output()
	lines := strings.Split(string(out), "\n")
	var s string = ""
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Index(line, "*") == 0 {
			s = line[2:len(line)]
		}
	}
	return s
}

func isWorkingTreeDirty() bool {
	out, _ := exec.Command("git", "status").Output()
	return strings.Index(string(out), "nothing to commit, working tree clean") == -1
}
