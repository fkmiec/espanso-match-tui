package main

import (
	"encoding/json"
	"fmt"
)

/*
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
		Trigger string `json:"trigger,omitempty"`
		Regex   string `json:"regex,omitempty"`

		//Action to take in response to the match: simple replace, markdown or HTML replace, or populate an image
		Replace  string `json:"replace,omitempty"`
		Markdown string `json:"markdown,omitempty"`
		Html     string `json:"html,omitempty"`
		Image    string `json:"image_path,omitempty"`

		//Espanso's defined attributes for a match.
		Label          string `json:"label,omitempty"`
		SearchTerms    string `json:"search_terms,omitempty"`
		Word           bool   `json:"word,omitempty"`
		PropagateCase  bool   `json:"propogate_case,omitempty"`
		UppercaseStyle string `json:"uppercase_style,omitempty"`
		Paragraph      bool   `json:"paragraph,omitempty"` //Only relevant for markdown

		Anchor string `json:"anchor,omitempty"`
		Vars   []Var  `json:"vars,omitempty"`
	}

// Alias to avoid recursive call
type AliasMatch Match
*/
func (m *Match) UnmarshalJSON(bytes []byte) error {
	fmt.Printf("Attempting Match.UnmarshalJSON() for %s\n", string(bytes))

	var alias AliasMatch

	if err := json.Unmarshal(bytes, &alias); err != nil {
		fmt.Printf("Error in Match.UnmarshalJSON(). Error: %s\nData: %s\n", err.Error(), string(bytes))
		return err
	}

	*m = Match(alias)

	fmt.Printf("Match.UnmarshalJSON: After json.Unmarshal() for %s\n", string(bytes))

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

/*
	type Var struct {
		Name string `json:"name"`
		Type string `json:"type"`
		//ParamsJSON json.RawMessage `json:"params,omitempty"`
		Params any `json:"params,omitempty"`
		//Implement a custom decoder that will instantiate the correct params for the variable type
		//NOTE: type = clipboard has only name and type, so use simple Var for that one
	}

	type DateParams struct {
		Format string `json:"format"`
		Locale string `json:"locale,omitempty"`
		Offset int    `json:"offset,omitempty"`
	}

	type EchoParams struct {
		Echo string `json:"echo"`
	}

	type TriggerParams struct {
		Trigger string `json:"trigger"`
	}

	type Choice struct {
		Label string `json:"label"`
		Value string `json:"id,omitempty"`
	}

	type ChoiceParams struct {
		Values []Choice `json:"values"`
	}

	type RandomParams struct {
		Choices []string `json:"choices"`
	}

	type ScriptParams struct {
		Args []string `json:"args"`
	}

	type ShellParams struct {
		Cmd   string `json:"cmd"`
		Shell string `json:"shell,omitempty"`
		Trim  bool   `json:"trim,omitempty"`
		Debug bool   `json:"debug,omitempty"`
	}

	type Field struct {
		Multiline bool     `json:"multiline"`
		Type      string   `json:"type"`
		Values    []Choice `json:"values"`
	}

	type FormParams struct {
		Layout string           `json:"layout"`
		Fields map[string]Field `json:"fields"`
	}
*/
type AliasJSONVar struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	ParamsJSON json.RawMessage `json:"params,omitempty"`
	Params     any
}

func (v *Var) UnmarshalJSON(bytes []byte) error {

	var alias AliasJSONVar
	if err := json.Unmarshal(bytes, &alias); err != nil {
		return err
	}
	v.Name = alias.Name
	v.Type = alias.Type

	switch alias.Type {
	case "date":
		var params DateParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "echo":
		var params EchoParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "trigger":
		var params TriggerParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "choice":
		var params ChoiceParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "random":
		var params RandomParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "script":
		var params ScriptParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "shell":
		var params ShellParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
			return err
		}
		v.Params = params
	case "form":
		var params FormParams
		if err := json.Unmarshal(alias.ParamsJSON, &params); err != nil {
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
type MatchFile struct {
	Matches []Match `json:"matches"`
}
*/
/*
func (f *MatchFile) UnmarshalJSON(bytes []byte) error {
	//Unmarshall to map.
	// Get list of matches.
	// For each one, pass the raw data for the match to the Match struct's UnmarshalJSON method
	var m map[string][]json.RawMessage
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}
	list, ok := m["matches"]
	if !ok {
		panic("Did not find list of matches in the file! Exiting.")
	}
	var match Match
	for _, raw := range list {
		match = Match{}
		if err := match.UnmarshalJSON(raw); err != nil {
			errMsg := fmt.Sprintf("Failed to unmarshal Match. Error: %s\nData: %s", err.Error(), string(raw))
			panic(errMsg)
		}
		f.Matches = append(f.Matches, match)
	}
	return nil
}
*/

func main() {
	fmt.Println("Hello World!")
	jsonString := `{
    "matches": [
        {
            "trigger": ":espanso",
            "replace": "Hi there!"
        },
        {
            "trigger": ":date",
            "replace": "{{mydate}}",
            "label": "Date",
            "vars": [
                {
                    "name": "mydate",
                    "type": "date",
                    "params": {
                        "format": "%m/%d/%Y"
                    }
                }
            ]
        },
        {
            "trigger": ":wiki",
            "replace": "{{output}}",
            "label": "Open A Tiddlywiki",
            "vars": [
                {
                    "name": "output",
                    "type": "shell",
                    "params": {
                        "cmd": "|\ndotool <<< 'key BackSpace BackSpace BackSpace BackSpace BackSpace'\ngoogle-chrome file:///home/fkmiec/Downloads/tiddlywiki.html"
                    }
                }
            ]
        }
    ]
}
	`
	file := MatchFile{}
	if err := json.Unmarshal([]byte(jsonString), &file); err != nil {
		panic("Marshaling error: " + err.Error())
	}
	//if err := file.UnmarshalJSON([]byte(jsonString)); err != nil {
	//	panic("Marshaling error: " + err.Error())
	//}

	for _, m := range file.Matches {
		fmt.Printf("Match Trigger:\n%s\n", m.Trigger)
		fmt.Printf("Match Vars:\n%+v\n", m.Vars)
	}

	fmt.Println("Printing re-marshaled JSON")
	afterMarshal, err := json.Marshal(file)
	if err != nil {
		fmt.Printf("Marshal Error: %s\n", err.Error())
	} else {
		fmt.Println(string(afterMarshal))
	}

}
