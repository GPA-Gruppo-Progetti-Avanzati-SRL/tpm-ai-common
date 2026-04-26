package anthropiclks

import (
	"context"
	"errors"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/store/promptqueueitem"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/mongolks"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ParseMessage(p prompts.PromptTemplate, message *anthropic.Message) (map[string]prompts.MessagePart, error) {
	const semLogContext = "parse-message"
	var err error

	if len(message.Content) != 1 {
		err = errors.New("expected a single message")
		log.Warn().Err(err).Msg(semLogContext)
		return nil, err
	}

	for _, m := range message.Content {
		if m.Type == "text" {
			return p.ParseTextContent(m.Text)
		}
	}

	return nil, nil
}

func FindPromptQueueItemByObjectId(objectId string) (*promptqueueitem.PromptQueueItem, error) {
	const semLogContext = semLogContextBase + "retrieve-item-from-object-id"

	coll, err := mongolks.GetCollection(context.Background(), "default", promptqueueitem.CollectionId)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	item, _, err := promptqueueitem.FindByObjectId(coll, objectId, true, options.FindOne())
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return item, nil
}
