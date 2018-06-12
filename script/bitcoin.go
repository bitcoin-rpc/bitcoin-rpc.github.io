package main

import (
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Name      string
	Groupname string
}

type Group struct {
	Name     string
	Commands []Command
}

type Document struct {
	Command *Command
	Groups  []Group
}

func main() {
	first := run("help")
	split := strings.Split(first, "\n")

	groups := make([]Group, 0)
	commands := make([]Command, 0)
	lastGroupName := ""

	for _, line := range split {
		if len(line) > 0 {
			if strings.HasPrefix(line, "== ") {
				if len(commands) != 0 {
					g := Group{
						Name:     lastGroupName,
						Commands: commands,
					}
					groups = append(groups, g)
					commands = make([]Command, 0)
				}
				lastGroupName = line[3 : len(line)-3]
			} else {
				name := strings.Split(line, " ")[0]
				comm := Command{
					Name:      name,
					Groupname: strings.ToLower(lastGroupName),
				}
				commands = append(commands, comm)
			}
		}
	}

	g := Group{
		Name:     lastGroupName,
		Commands: commands,
	}
	groups = append(groups, g)

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(open("../index.html"), Document{
		Command: nil,
		Groups:  groups,
	})

	for _, group := range groups {
		for _, command := range group.Commands {
			tmpl.Execute(open("../"+command.Name+".html"), Document{
				Command: &command,
				Groups:  groups,
			})
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func open(path string) io.Writer {
	f, err := os.Create(path)
	// not closing, program will close sooner
	check(err)
	return f
}

func run(args ...string) string {
	out, err := exec.Command("/home/g/dev/bitcoind/bitcoin-0.16.0/bin/bitcoin-cli", args...).CombinedOutput()
	if err != nil {
		panic(err)
	}

	return string(out)
}
