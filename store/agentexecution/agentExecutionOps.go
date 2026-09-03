package agentexecution

import (
	"context"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/cob-game/store/commons"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/mongolks"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewAgentExecution(domain, site string, agentBid, jobId string, bidRef commons.BidTextPair, status string, weight int32, params bson.M) (*AgentExecution, error) {
	const semLogContext = semLogPackageContext + "new-agent-execution"
	itemColl, err := mongolks.GetCollection(context.Background(), "default", CollectionId)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	item := AgentExecution{
		Domain: domain,
		Site:   site,
		Bid:    agentBid,
		Et:     EntityType,
		Status: status,
		BidRef: BidEtPair{
			Bid: bidRef.Bid,
			Et:  bidRef.Et,
		},
		Weight: weight,
		Params: params,
		Group:  jobId,
	}

	resp, err := itemColl.InsertOne(context.Background(), item)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	log.Info().Interface("resp", resp).Msg(semLogContext)
	item.OId = resp.InsertedID.(bson.ObjectID)
	return &item, nil
}
