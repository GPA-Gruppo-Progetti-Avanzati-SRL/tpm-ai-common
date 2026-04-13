package prompts

type Category struct {
	Name    string                 `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Lks     string                 `yaml:"lks,omitempty" mapstructure:"lks,omitempty" json:"lks,omitempty"`
	Model   string                 `yaml:"model,omitempty" mapstructure:"model,omitempty" json:"model,omitempty"`
	Prompt  string                 `yaml:"prompt,omitempty" mapstructure:"prompt,omitempty" json:"prompt,omitempty"`
	Options map[string]interface{} `yaml:"options,omitempty" mapstructure:"options,omitempty" json:"options,omitempty"`
}

func (cat Category) GetStringOption(n string, defValue string) string {

	if v, ok := cat.Options[n]; ok {
		if vs, ok := v.(string); ok {
			return vs
		}
	}

	return defValue
}

func (cat Category) GetFloat64Option(n string, defValue float64) float64 {

	if v, ok := cat.Options[n]; ok {
		if vf, ok := v.(float64); ok {
			return vf
		}
	}

	return defValue
}

func (cat Category) GetIntOption(n string, defValue int) int {

	if v, ok := cat.Options[n]; ok {
		if vi, ok := v.(int); ok {
			return vi
		}
	}

	return defValue
}

func (cat Category) GetInt64Option(n string, defValue int64) int64 {

	if v, ok := cat.Options[n]; ok {
		if vi, ok := v.(int); ok {
			return int64(vi)
		}
	}

	return defValue
}
