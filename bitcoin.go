package main

import (
	"html"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Command struct {
	name        string
	description string
}

type Group struct {
	name     string
	commands []Command
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
						name:     lastGroupName,
						commands: commands,
					}
					groups = append(groups, g)
					commands = make([]Command, 0)
				}
				lastGroupName = line[3 : len(line)-3]
			} else {
				name := strings.Split(line, " ")[0]
				desc := run("help", name)
				desc = html.EscapeString(desc)
				comm := Command{
					name:        name,
					description: desc,
				}
				commands = append(commands, comm)
			}
		}
	}

	g := Group{
		name:     lastGroupName,
		commands: commands,
	}
	groups = append(groups, g)

	menuStr := ""
	for _, group := range groups {
		menuStr += `
<div class="card">
  <div class="card-header">
    `
		menuStr += group.name
		menuStr += `
  </div>
  <div class="card-body">`

		for _, command := range group.commands {
			menuStr += `
    <a href="`
			menuStr += command.name
			menuStr += `.html">`
			menuStr += command.name
			menuStr += `</a><br>`
		}

		menuStr += `
  </div>
</div>
<br>`

	}
	check(ioutil.WriteFile("tmp-menu.inc", []byte(menuStr), 0644))

	runBash("cat start.inc index.inc mid.inc tmp-menu.inc end.inc | sed 's/XXX/Main page/' > index.html")

	for _, group := range groups {
		for _, command := range group.commands {
			check(ioutil.WriteFile("tmp.inc", []byte(command.description), 0644))
			runBash("cat start.inc tmp.inc mid.inc tmp-menu.inc end.inc | sed 's/XXX/" + command.name + "/' > " + command.name + ".html")
		}
	}
	runBash("rm tmp.inc")
	runBash("rm tmp-menu.inc")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runBash(cmd string) {
	check(exec.Command("bash", "-c", cmd).Run())
}

func run(args ...string) string {
	out, err := exec.Command("/home/g/dev/bitcoind/bitcoin-0.16.0/bin/bitcoin-cli", args...).CombinedOutput()
	if err != nil {
		panic(err)
	}

	return string(out)
}
