package promptqueueitem

import (
	"context"
	"errors"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func UpdateStatusByObjectId(coll *mongo.Collection, objectId string, status string, batchId string) error {
	const semLogContext = semLogPackageContext + "update-status-by-object-id"
	obId, err := bson.ObjectIDFromHex(objectId)
	if err != nil {
		return err
	}

	updOpts := []UpdateOption{
		UpdateWithStatus(status),
	}

	if batchId != "" {
		updOpts = append(updOpts, UpdateWithBatch_id(batchId))
	}

	filter := Filter{}
	filter.Or().AndOIdEqTo(obId)

	fd := filter.Build()
	updDoc := GetUpdateDocumentFromOptions(updOpts...)
	ud := updDoc.Build()
	log.Trace().
		Str("filter", util.MustToExtendedJsonString(fd, false, false)).
		Str("update", util.MustToExtendedJsonString(ud, false, false)).
		Msg(semLogContext)

	resp, err := coll.UpdateOne(context.Background(), fd, ud)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	if resp.MatchedCount == 0 {
		err = errors.New("no document updated matched")
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	return nil
}

func UpdateStatusByListOfObjectId(coll *mongo.Collection, objectId []string, status string, batchId string) error {
	const semLogContext = semLogPackageContext + "update-status-by-list-of-object-id"

	var objIds []bson.ObjectID
	for _, oid := range objectId {
		obId, err := bson.ObjectIDFromHex(oid)
		if err != nil {
			return err
		}

		objIds = append(objIds, obId)
	}

	updOpts := []UpdateOption{
		UpdateWithStatus(status),
		UpdateWithBatch_id(batchId),
	}

	filter := Filter{}
	filter.Or().AndOIdIn(objIds)

	fd := filter.Build()
	updDoc := GetUpdateDocumentFromOptions(updOpts...)
	ud := updDoc.Build()
	log.Trace().
		Str("filter", util.MustToExtendedJsonString(fd, false, false)).
		Str("update", util.MustToExtendedJsonString(ud, false, false)).
		Msg(semLogContext)

	resp, err := coll.UpdateOne(context.Background(), fd, ud)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	if resp.MatchedCount == 0 {
		err = errors.New("no document updated matched")
		log.Error().Err(err).Msg(semLogContext)
		return err
	}

	return nil
}
