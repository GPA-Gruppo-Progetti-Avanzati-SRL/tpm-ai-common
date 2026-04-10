package anthropiclks

import (
	"errors"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-ai-common/linkedservices/prompts"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog/log"
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
