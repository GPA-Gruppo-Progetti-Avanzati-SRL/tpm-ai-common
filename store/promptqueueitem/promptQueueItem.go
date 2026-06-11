package promptqueueitem

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// @tpm-schematics:start-region("top-file-section")
const (
	CollectionId         = "prompt-queue-item"
	EntityType           = "prompt-queue-item"
	semLogPackageContext = "propmt-queue-item::"

	StatusReady      = "ready"
	StatusProcessing = "processing"
	StatusFailed     = "failed"
	StatusCompleted  = "completed"
)

type Counter struct {
	Dimension1 string `json:"dimension_1,omitempty" bson:"dimension_1,omitempty" yaml:"dimension_1,omitempty"`
	Count1     int32  `json:"count_1,omitempty" bson:"count_1,omitempty" yaml:"count_1,omitempty"`
	Count2     int32  `json:"count_2,omitempty" bson:"count_2,omitempty" yaml:"count_2,omitempty"`
}

type CounterQueryResult struct {
	Records int       `json:"records,omitempty" bson:"records,omitempty" yaml:"records,omitempty"`
	Data    []Counter `json:"data,omitempty" bson:"data,omitempty" yaml:"data,omitempty"`
}

// @tpm-schematics:end-region("top-file-section")

type PromptQueueItem struct {
	OId        bson.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty" yaml:"_id,omitempty"`
	Domain     string         `json:"domain,omitempty" bson:"domain,omitempty" yaml:"domain,omitempty"`
	Site       string         `json:"site,omitempty" bson:"site,omitempty" yaml:"site,omitempty"`
	Bid        string         `json:"_bid,omitempty" bson:"_bid,omitempty" yaml:"_bid,omitempty"`
	Et         string         `json:"_et,omitempty" bson:"_et,omitempty" yaml:"_et,omitempty"`
	Category   string         `json:"category,omitempty" bson:"category,omitempty" yaml:"category,omitempty"`
	Status     string         `json:"status,omitempty" bson:"status,omitempty" yaml:"status,omitempty"`
	BatchMode  bool           `json:"batch_mode,omitempty" bson:"batch_mode,omitempty" yaml:"batch_mode,omitempty"`
	BatchId    string         `json:"batch_id,omitempty" bson:"batch_id,omitempty" yaml:"batch_id,omitempty"`
	BidRef     BidEtPair      `json:"bid_ref,omitempty" bson:"bid_ref,omitempty" yaml:"bid_ref,omitempty"`
	Weight     int32          `json:"weight,omitempty" bson:"weight,omitempty" yaml:"weight,omitempty"`
	BucketPath BucketPathPair `json:"bucketPath,omitempty" bson:"bucketPath,omitempty" yaml:"bucketPath,omitempty"`
	Metadata   bson.M         `json:"metadata,omitempty" bson:"metadata,omitempty" yaml:"metadata,omitempty"`
	Singleton  bool           `json:"singleton,omitempty" bson:"singleton,omitempty" yaml:"singleton,omitempty"`

	// @tpm-schematics:start-region("struct-section")
	Count int32 `json:"count,omitempty" bson:"count,omitempty" yaml:"count,omitempty"`
	// @tpm-schematics:end-region("struct-section")
}

func (s PromptQueueItem) IsZero() bool {
	return s.OId == bson.NilObjectID && s.Domain == "" && s.Site == "" && s.Bid == "" && s.Et == "" && s.Category == "" && s.Status == "" && !s.BatchMode && s.BatchId == "" && s.BidRef.IsZero() && s.Weight == 0 && s.BucketPath.IsZero() && len(s.Metadata) == 0 && !s.Singleton
}

type QueryResult struct {
	Records int               `json:"records,omitempty" bson:"records,omitempty" yaml:"records,omitempty"`
	Data    []PromptQueueItem `json:"data,omitempty" bson:"data,omitempty" yaml:"data,omitempty"`
}

type FormResponseError struct {
	Field string `json:"field,omitempty" bson:"field,omitempty" yaml:"field,omitempty"`
	Error string `json:"message,omitempty" bson:"message,omitempty" yaml:"message,omitempty"`
}

type FormResponse struct {
	Status      int                 `json:"status,omitempty" bson:"status,omitempty" yaml:"status,omitempty"`
	Message     string              `json:"message,omitempty" bson:"message,omitempty" yaml:"message,omitempty"`
	FieldErrors []FormResponseError `json:"fieldErrors,omitempty" bson:"fieldErrors,omitempty" yaml:"fieldErrors,omitempty"`
	Document    *PromptQueueItem    `json:"document,omitempty" bson:"document,omitempty" yaml:"document,omitempty"`
}

// @tpm-schematics:start-region("bottom-file-section")

type PromptQueueItems []PromptQueueItem

func (items PromptQueueItems) ToListOfObjectIds() []string {
	var objectIds []string
	for _, e := range items {
		objectIds = append(objectIds, e.OId.Hex())
	}

	return objectIds
}

func (items PromptQueueItems) ToCSVListOfObjectIds() string {
	return strings.Join(items.ToListOfObjectIds(), ",")
}

func (ce PromptQueueItem) Log(logZeroEvt *zerolog.Event) *zerolog.Event {
	logZeroEvt.Str("category", ce.Category).Str("status", ce.Status).Str("batch-id", ce.BatchId).Str("bid", ce.BidRef.Bid).Int32("weight", ce.Weight).Str("bucket", ce.BucketPath.String())
	return logZeroEvt
}

func (ce PromptQueueItem) ParseMetadata() (map[string]interface{}, error) {
	const semLogContext = semLogPackageContext + "parse-metadata"
	if len(ce.Metadata) == 0 {
		return nil, nil
	}

	m := make(map[string]interface{})
	for k, v := range ce.Metadata {
		switch vt := v.(type) {
		case string:
			m[k] = vt
		case bson.A:
			for ivv, vv := range vt {
				switch vvt := vv.(type) {
				case string:
					if ivv == 0 {
						m[k] = []string{}
					}
					m[k] = append(m[k].([]string), vvt)
				case bson.D:
					switch len(vvt) {
					case 2:
						if vvt[0].Key == "bucket" && vvt[1].Key == "path" {
							m[k] = BucketPathPair{Bucket: fmt.Sprint(vvt[0].Value), Path: fmt.Sprint(vvt[1].Value)}
						}
					default:
						err := errors.New("invalid metadata")
						log.Error().Err(err).Msg(semLogContext)
						return nil, err
					}
				}
			}
		}
	}

	return m, nil
}

// @tpm-schematics:end-region("bottom-file-section")
