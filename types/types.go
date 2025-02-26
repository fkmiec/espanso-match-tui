package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type MatchFile struct {
	Matches []Match `yaml:"matches"`
}

type MatchType int

const (
	TRIGGER_REPLACE MatchType = iota
	TRIGGER_MARKDOWN
	TRIGGER_HTML
	TRIGGER_IMAGE
	REGEX_REPLACE
	REGEX_MARKDOWN
	REGEX_HTML
	REGEX_IMAGE
)

type Match struct {
	//Helper variable to be populated in custom UnmarshalJSON
	Type MatchType

	//Ordinary or Regex trigger for the match
	Trigger string `yaml:"trigger,omitempty"`
	Regex   string `yaml:"regex,omitempty"`

	//Action to take in response to the match: simple replace, markdown or HTML replace, or populate an image
	Replace  string `yaml:"replace,omitempty"`
	Markdown string `yaml:"markdown,omitempty"`
	Html     string `yaml:"html,omitempty"`
	Image    string `yaml:"image_path,omitempty"`

	//Espanso's defined attributes for a match.
	Label          string `yaml:"label,omitempty"`
	SearchTerms    string `yaml:"search_terms,omitempty"`
	Word           bool   `yaml:"word,omitempty"`
	PropagateCase  bool   `yaml:"propogate_case,omitempty"`
	UppercaseStyle string `yaml:"uppercase_style,omitempty"`
	Paragraph      bool   `yaml:"paragraph,omitempty"` //Only relevant for markdown

	Anchor string `yaml:"anchor,omitempty"`
	Vars   []Var  `yaml:"vars,omitempty"`
}

// Alias to avoid recursive call
type AliasMatch Match

// UnmarshalYAML(bytes []byte) error {
func (m *Match) UnmarshalYAML(value *yaml.Node) error {
	fmt.Printf("Attempting Match.UnmarshalYAML() for %s\n", value.Value)

	var alias AliasMatch

	if err := value.Decode(&alias); err != nil {
		fmt.Printf("Error in Match.UnmarshalYAML(). Error: %s\nData: %s\n", err.Error(), value.Value)
		return err
	}

	*m = Match(alias)

	fmt.Printf("Match.UnmarshalYAML: After yaml.Unmarshal() for %s\n", value.Value)

	if m.Trigger != "" {
		if m.Replace != "" {
			m.Type = TRIGGER_REPLACE
		} else if m.Markdown != "" {
			m.Type = TRIGGER_MARKDOWN
		} else if m.Html != "" {
			m.Type = TRIGGER_HTML
		} else if m.Image != "" {
			m.Type = TRIGGER_IMAGE
		}
	} else if m.Regex != "" {
		if m.Replace != "" {
			m.Type = REGEX_REPLACE
		} else if m.Markdown != "" {
			m.Type = REGEX_MARKDOWN
		} else if m.Html != "" {
			m.Type = REGEX_HTML
		} else if m.Image != "" {
			m.Type = REGEX_IMAGE
		}
	}
	return nil
}

type Var struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	//ParamsYAML yaml.RawMessage `yaml:"params,omitempty"`
	Params any `yaml:"params,omitempty"`
	//Implement a custom decoder that will instantiate the correct params for the variable type
	//NOTE: type = clipboard has only name and type, so use simple Var for that one
}

type DateParams struct {
	Format string `yaml:"format"`
	Locale string `yaml:"locale,omitempty"`
	Offset int    `yaml:"offset,omitempty"`
}

type EchoParams struct {
	Echo string `yaml:"echo"`
}

type TriggerParams struct {
	Trigger string `yaml:"trigger"`
}

type Choice struct {
	Label string `yaml:"label"`
	Value string `yaml:"id,omitempty"`
}

type ChoiceParams struct {
	Values []Choice `yaml:"values"`
}

type RandomParams struct {
	Choices []string `yaml:"choices"`
}

type ScriptParams struct {
	Args []string `yaml:"args"`
}

type ShellParams struct {
	Cmd   string `yaml:"cmd"`
	Shell string `yaml:"shell,omitempty"`
	Trim  bool   `yaml:"trim,omitempty"`
	Debug bool   `yaml:"debug,omitempty"`
}

type Field struct {
	Multiline bool     `yaml:"multiline"`
	Type      string   `yaml:"type"`
	Values    []Choice `yaml:"values"`
}

type FormParams struct {
	Layout string           `yaml:"layout"`
	Fields map[string]Field `yaml:"fields"`
}

type AliasVar struct {
	Name       string    `yaml:"name"`
	Type       string    `yaml:"type"`
	ParamsYAML yaml.Node `yaml:"params,omitempty"`
}

func (v *Var) UnmarshalYAML(value *yaml.Node) error {

	var alias AliasVar
	if err := value.Decode(&alias); err != nil {
		return err
	}
	v.Name = alias.Name
	v.Type = alias.Type

	switch alias.Type {
	case "date":
		var params DateParams
		//if err := yaml.Unmarshal(alias.ParamsYAML, &params); err != nil {
		//	return err
		//}
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "echo":
		var params EchoParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "trigger":
		var params TriggerParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "choice":
		var params ChoiceParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "random":
		var params RandomParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "script":
		var params ScriptParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "shell":
		var params ShellParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	case "form":
		var params FormParams
		if err := alias.ParamsYAML.Decode(&params); err != nil {
			return err
		}
		v.Params = params
	}
	//v.Params = alias.Params

	//fmt.Printf("Var.UnmarshalJSON case 'date' (orig): %+v\n", v.Params)
	//fmt.Printf("Var.UnmarshalJSON case 'date' (parsed): %+v\n", v)

	return nil
}

/*
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

	file := MatchFile{}
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
*/
