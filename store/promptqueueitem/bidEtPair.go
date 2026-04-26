package promptqueueitem

import "fmt"

// @tpm-schematics:start-region("top-file-section")
// @tpm-schematics:end-region("top-file-section")

type BidEtPair struct {
	Bid string `json:"bid,omitempty" bson:"bid,omitempty" yaml:"bid,omitempty"`
	Et  string `json:"et,omitempty" bson:"et,omitempty" yaml:"et,omitempty"`

	// @tpm-schematics:start-region("struct-section")
	// @tpm-schematics:end-region("struct-section")
}

func (s BidEtPair) IsZero() bool {
	return s.Bid == "" && s.Et == ""
}

// @tpm-schematics:start-region("bottom-file-section")

func (bidPair BidEtPair) String() string {
	return fmt.Sprintf("%s-%s", bidPair.Et, bidPair.Bid)

}

// @tpm-schematics:end-region("bottom-file-section")
