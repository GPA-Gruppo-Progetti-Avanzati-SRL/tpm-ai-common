package promptqueueitem

import (
	"context"
	"fmt"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	BidTaskPropertyName            = "bid"
	BidsTaskPropertyName           = "bid"
	EtTaskPropertyName             = "et"
	CategoryTaskPropertyName       = "category"
	BatchAPITaskPropertyName       = "batch-api"
	DestBlobBucketTaskPropertyName = "dest-blob-bucket"
	DestBlobPathTaskPropertyName   = "dest-blob-path"
)

func PushPromptQueueItem(coll *mongo.Collection, item *PromptQueueItem) error {
	const semLogContext = semLogPackageContext + "push-prompt-queue-item"

	item.Bid = fmt.Sprintf("%s-%s-%s", item.BidRef.Et, item.Category, util.NewUUID())
	item.Status = StatusReady

	_, err := coll.InsertOne(context.Background(), item)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	return nil
}

/*


 */

/*
func paramsFromProperties(properties map[string]interface{}) (PromptQueueItem, error) {
	const semLogContext = semLogPackageContext + "params-from-properties"
	var params PromptQueueItem

	bidsRaw := getStringProperty(properties, BidsTaskPropertyName)
	if bidsRaw != "" {
		for _, b := range strings.Split(bidsRaw, ",") {
			if b = strings.TrimSpace(b); b != "" {
				params.bids = append(params.bids, b)
			}
		}
	}
	if len(params.bids) == 0 {
		if singleBid := aTask.GetStringProperty(BidTaskPropertyName, task.PropertiesTaskAndPartitionScope); singleBid != "" {
			params.bids = []string{singleBid}
		}
	}
	if len(params.bids) == 0 {
		w.wrkLogger.Logger.Info().Msg(semLogContext + " No bid(s) found in task properties")
		return params, errors.New("missing program bid(s)")
	}

	params.category = aTask.GetStringProperty(CategoryTaskPropertyName, task.PropertiesTaskAndPartitionScope)
	if params.category == "" {
		w.wrkLogger.Logger.Info().Msg(semLogContext + " No category found in task properties")
		return params, errors.New("missing category")
	}

	params.et = aTask.GetStringProperty(EtTaskPropertyName, task.PropertiesTaskAndPartitionScope)
	if params.et == "" {
		w.wrkLogger.Logger.Info().Msg(semLogContext + " No entity-type found in task properties")
		return params, errors.New("missing entity-type")
	}

	params.destBlobBucket = aTask.GetStringProperty(DestBlobBucketTaskPropertyName, task.PropertiesTaskAndPartitionScope)
	if params.et == "" {
		w.wrkLogger.Logger.Info().Msg(semLogContext + " No dest-blob-bucket found in task properties")
		return params, errors.New("no dest-blob-bucket found")
	}

	params.destBlobPath = aTask.GetStringProperty(DestBlobPathTaskPropertyName, task.PropertiesTaskAndPartitionScope)
	if params.destBlobPath == "" {
		w.wrkLogger.Logger.Info().Msg(semLogContext + " No dest-blob-path found in task properties")
		return params, errors.New("no dest-blob-path found")
	}


	return params, nil
}
*/

func getStringProperty(properties map[string]interface{}, key string) string {
	if v, ok := properties[key]; ok {
		if sv, ok := v.(string); ok {
			return sv
		}
	}

	return ""
}

func getBoolProperty(properties map[string]interface{}, key string) (bool, bool) {
	if v, ok := properties[key]; ok {
		if sv, ok := v.(bool); ok {
			return sv, true
		}
	}

	return false, false
}
