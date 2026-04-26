package promptqueueitem

import "fmt"

// @tpm-schematics:start-region("top-file-section")
// @tpm-schematics:end-region("top-file-section")

type BucketPathPair struct {
	Bucket string `json:"bucket,omitempty" bson:"bucket,omitempty" yaml:"bucket,omitempty"`
	Path   string `json:"path,omitempty" bson:"path,omitempty" yaml:"path,omitempty"`

	// @tpm-schematics:start-region("struct-section")
	// @tpm-schematics:end-region("struct-section")
}

func (s BucketPathPair) IsZero() bool {
	return s.Bucket == "" && s.Path == ""
}

// @tpm-schematics:start-region("bottom-file-section")

func (s BucketPathPair) String() string {
	return fmt.Sprintf("Bucket: %s, Path: %s", s.Bucket, s.Path)
}

// @tpm-schematics:end-region("bottom-file-section")
