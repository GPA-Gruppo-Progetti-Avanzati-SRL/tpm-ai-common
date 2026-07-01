package promptqueueitem

import (
	"context"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Delete(coll *mongo.Collection, filter *Filter, limit int) (int, error) {
	const semLogContext = semLogPackageContext + "delete"

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	qr, err := Find(coll, filter, false, opts)
	if err != nil {
		return 0, err
	}

	if len(qr.Data) == 0 {
		return 0, nil
	}

	var totalDeletions int
	for i := 0; i < len(qr.Data); i += 20 {

		var ids []string
		for ndx := i; ndx < i+20 && ndx < len(qr.Data); ndx++ {
			ids = append(ids, qr.Data[ndx].Bid)
		}

		f1 := Filter{}
		f1.Or().AndBidIn(ids).AndEtEqTo(EntityType)
		f1d := f1.Build()
		log.Info().Str("filter", util.MustToExtendedJsonString(f1d, false, false)).Msg(semLogContext + " - deleting documents")
		resp, err := coll.DeleteMany(context.Background(), f1d, options.DeleteMany())
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return 0, err
		}

		totalDeletions += int(resp.DeletedCount)
	}

	return totalDeletions, nil
}
