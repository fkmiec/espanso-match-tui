package main

import (
	types "espanso-match-tui/types"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("Hello World!")
	yamlString := `
# espanso match file

# For a complete introduction, visit the official docs at: https://espanso.org/docs/

# You can use this file to define the base matches (aka snippets)
# that will be available in every application when using espanso.

# Matches are substitution rules: when you type the "trigger" string
# it gets replaced by the "replace" string.
matches:
  # Simple text replacement
  - trigger: ":espanso"
    replace: "Hi there!"

  # NOTE: espanso uses YAML to define matches, so pay attention to the indentation!

  # But matches can also be dynamic:

  # Print the current date
  - trigger: ":date"
    replace: "{{mydate}}"
    label: "Date"
    vars:
      - name: mydate
        type: date
        params:
          format: "%m/%d/%Y"

  # Open browser to a specified url
  - trigger: ":wiki"
    replace: "{{output}}" 
    label: "Open A Tiddlywiki"
    vars:
      - name: output
        type: shell
        params:
          cmd: |
            dotool <<< 'key BackSpace BackSpace BackSpace BackSpace BackSpace'
            google-chrome file:///home/fkmiec/Downloads/tiddlywiki.html

  # Exec python script
  - trigger: ":py"
    replace: "{{output}}"
    label: "Python Clip Script"
    vars:
      - name: output
        type: shell
        params:
          cmd: "python3 /home/fkmiec/.config/espanso/scripts/clip.py"

  # Exec python script
  - trigger: ":pr"
    replace: "{{output}}"
    label: "Python Version"
    vars:
      - name: output
        type: shell
        params:
          cmd: "python3 --version"
`

	file := types.MatchFile{}
	if err := yaml.Unmarshal([]byte(yamlString), &file); err != nil {
		panic("Marshaling error: " + err.Error())
	}

	for _, m := range file.Matches {
		fmt.Printf("Match Trigger:\n%s\n", m.Trigger)
		fmt.Printf("Match Vars:\n%+v\n", m.Vars)
	}

	fmt.Println("Printing re-marshaled YAML")
	afterMarshal, err := yaml.Marshal(file)
	if err != nil {
		fmt.Printf("Marshal Error: %s\n", err.Error())
	} else {
		fmt.Println(string(afterMarshal))
	}

}
