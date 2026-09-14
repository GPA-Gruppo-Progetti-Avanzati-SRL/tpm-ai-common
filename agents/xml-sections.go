package agents

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// XMLSection describes a single named part of an XMLSectionedDocument. The layout
// mirrors the YAML representation found in xml-document-sections.yml: each entry
// declares how a delimited section of an LLM response should be interpreted and
// materialized (content type, extension, whether it is a standalone document,
// etc.).
type XMLSection struct {
	Name        string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Ct          string `yaml:"ct,omitempty" mapstructure:"ct,omitempty" json:"ct,omitempty"`
	Ext         string `yaml:"ext,omitempty" mapstructure:"ext,omitempty" json:"ext,omitempty"`
	IsDocument  bool   `yaml:"is-document,omitempty" mapstructure:"is-document,omitempty" json:"is-document,omitempty"`
	Required    bool   `yaml:"required,omitempty" mapstructure:"required,omitempty" json:"required,omitempty"`
	Title       string `yaml:"title,omitempty" mapstructure:"title,omitempty" json:"title,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Data        []byte `yaml:"-" mapstructure:"-" json:"-"`
}

type XMLSections []XMLSection

// NewXMLSectionsFromYAML unmarshals the YAML representation of a
// sections in a xml-sectioned document. The expected input is a top-level list of section
// descriptors, e.g.:
//
//   - name: overview
//     ext: "md"
//     is-document: true
//     ct: text/markdown
//     title: Overview
//     description: Overview del programma
func NewXMLSectionsFromYAML(b []byte) (XMLSections, error) {
	const semLogContext = semLogPackageContext + "new-xml-sections-from-yaml"

	var sections []XMLSection
	if err := yaml.Unmarshal(b, &sections); err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return sections, err
	}

	return sections, nil
}

// DelimiterNames returns the ordered list of section names, which double as the
// XML delimiters used to carve sections out of a text payload.
func (xs XMLSections) DelimiterNames() []string {
	names := make([]string, 0, len(xs))
	for _, s := range xs {
		names = append(names, s.Name)
	}
	return names
}

// SectionByName returns a pointer to the section with the given name, or nil if
// no such section exists.
func (xs XMLSections) SectionByName(name string) (XMLSection, bool) {
	for i := range xs {
		if xs[i].Name == name {
			return xs[i], true
		}
	}
	return XMLSection{}, false
}

// ExtractParts splits text into its XML-delimited sections using this
// document's section names as delimiters, returning a map keyed by section
// name. It does not mutate the document; use ExtractFromText to populate each
// section's Data instead.
func (xs XMLSections) ExtractParts(text string) (map[string]string, error) {
	return ExtractTextForParts(text, xs.DelimiterNames())
}

// ExtractFromText carves the XML-delimited sections out of text and populates
// each matching section's Data. A required section missing from the text yields
// an error.
func (xs XMLSections) ExtractFromText(text string) (XMLSections, error) {
	const semLogContext = semLogPackageContext + "extract-from-text"

	parts, err := xs.ExtractParts(text)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	var newXS XMLSections
	for _, sect := range xs {
		content, ok := parts[sect.Name]
		if !ok {
			if sect.Required {
				err = fmt.Errorf("required section %q not found in text", sect.Name)
				log.Error().Err(err).Msg(semLogContext)
				return nil, err
			}
			sect.Data = nil
		} else {
			sect.Data = []byte(content)
		}
		newXS = append(newXS, sect)
	}

	return newXS, nil
}
