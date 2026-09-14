package prompt

import (
	"path/filepath"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/agents"
)

// testPromptNames lists the prompt definitions under testdata/ to exercise. Each
// entry <name> maps to a testdata/<name>.yml definition file. Add new prompts
// here to have them covered by TestReadPromptDefinition.
var testPromptNames = []string{
	"node-summary",
}

func TestReadPromptDefinition(t *testing.T) {
	for _, name := range testPromptNames {
		t.Run(name, func(t *testing.T) {
			fn := filepath.Join("testdata", name+".yml")

			def, err := ReadPromptDefinition(fn)
			if err != nil {
				t.Fatalf("ReadPromptDefinition(%s) unexpected error: %v", fn, err)
			}

			if def.Name != name {
				t.Errorf("Name mismatch: got %q, want %q", def.Name, name)
			}

			// Location is set to the definition file's directory.
			if wantLoc := Location(filepath.Dir(fn)); def.Location != wantLoc {
				t.Errorf("Location mismatch: got %q, want %q", def.Location, wantLoc)
			}

			// ReadPromptDefinition already guarantees (on nil error) that SystemFn
			// and TemplateFn are resolved to existing files and that ParsedTemplate
			// is non-nil, so those are not re-asserted here.

			// UserPrompt exercises template execution, which ReadPromptDefinition
			// only parses but never runs.
			vars := make(map[string]string, len(def.Vars))
			for _, v := range def.Vars {
				vars[v.Name] = "./node-summary-source-code.cob"
			}
			if _, err := def.UserPrompt(vars, true, true); err != nil {
				t.Errorf("UserPrompt() unexpected error: %v", err)
			}

			// A definition with xml-sections and no schema outputs XML.
			if len(def.XMLSections) > 0 && def.OutputType() != agents.AgentResponseXML {
				t.Errorf("OutputType mismatch: got %q, want %q", def.OutputType(), agents.AgentResponseXML)
			}
		})
	}
}
