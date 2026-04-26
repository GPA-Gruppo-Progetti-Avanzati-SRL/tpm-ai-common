package promptqueueitem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-mongo-common/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// @tpm-schematics:start-region("top-file-section")
// @tpm-schematics:end-region("top-file-section")

// FindByPk ...
// @tpm-schematics:start-region("find-by-pk-signature-section")
func FindByObjectId(collection *mongo.Collection, objectId string, mustFind bool, findOptions *options.FindOneOptionsBuilder) (*PromptQueueItem, bool, error) {
	// @tpm-schematics:end-region("find-by-pk-signature-section")
	const semLogContext = "prompt-queue-item::find-by-pk"
	// @tpm-schematics:start-region("log-event-section")
	evtTraceLog := log.Trace()
	evtErrLog := log.Error()
	// @tpm-schematics:end-region("log-event-section")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ent := PromptQueueItem{}

	f := Filter{}
	// @tpm-schematics:start-region("filter-section")
	// customize the filtering
	oid, err := bson.ObjectIDFromHex(objectId)
	if err != nil {
		evtErrLog.Err(err).Msg(semLogContext)
		return nil, false, err
	}

	f.Or().AndOIdEqTo(oid)
	// @tpm-schematics:end-region("filter-section")
	fd := f.Build()
	evtTraceLog = evtTraceLog.Str("filter", util.MustToExtendedJsonString(fd, false, false))
	err = collection.FindOne(ctx, fd, findOptions).Decode(&ent)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		evtErrLog.Err(err).Msg(semLogContext)
		return nil, false, err
	} else {
		if err != nil {
			if mustFind {
				evtTraceLog.Msg(semLogContext + " document not found")
				return nil, false, err
			}

			evtTraceLog.Msg(semLogContext + " document not found but allowed")
			return nil, false, nil
		} else {
			evtTraceLog.Msg(semLogContext + " document found")
		}
	}

	return &ent, true, nil
}

func FindFirst(collection *mongo.Collection, f *Filter, findOptions *options.FindOptionsBuilder) (*PromptQueueItem, error) {
	const semLogContext = "prompt-queue-item::find-first"
	fd := f.Build()
	log.Trace().Str("filter", util.MustToExtendedJsonString(fd, false, false)).Msg(semLogContext)

	cur, err := collection.Find(context.Background(), fd, findOptions)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	if cur.Next(context.Background()) {
		dto := PromptQueueItem{}
		err = cur.Decode(&dto)
		if err != nil {
			return nil, err
		}

		return &dto, nil
	}

	return nil, nil
}

func Find(collection *mongo.Collection, f *Filter, withCount bool, findOptions *options.FindOptionsBuilder) (QueryResult, error) {
	const semLogContext = "prompt-queue-item::find"
	fd := f.Build()
	evtTraceLog := log.Trace().Str("filter", util.MustToExtendedJsonString(fd, false, false))
	evtErrLog := log.Error().Str("filter", util.MustToExtendedJsonString(fd, false, false))
	evtTraceLog.Msg(semLogContext)

	qr := QueryResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if withCount {
		countDocsOptions := options.Count()
		nr, err := collection.CountDocuments(ctx, fd, countDocsOptions)
		if err != nil {
			evtErrLog.Err(err).Msg(semLogContext)
			return qr, err
		}

		qr.Records = int(nr)
	}

	cur, err := collection.Find(ctx, fd, findOptions)
	if err != nil {
		evtErrLog.Err(err).Msg(semLogContext)
		return qr, err
	}

	for cur.Next(context.Background()) {
		dto := PromptQueueItem{}
		err = cur.Decode(&dto)
		if err != nil {
			return qr, err
		}

		qr.Data = append(qr.Data, dto)
	}

	if cur.Err() != nil {
		return qr, cur.Err()
	}

	return qr, nil
}

// @tpm-schematics:start-region("bottom-file-section")

type PromptsToProcessOptions struct {
	LimitCount         int
	LimitWeight        int
	MaxLimitWeight     int
	MaxStillProcessing int
}

func FindPromptsToProcess(coll *mongo.Collection, options PromptsToProcessOptions) ([]PromptQueueItem, bool, error) {
	const semLogContext = semLogPackageContext + "find-prompts-to-process"
	prompts, numProcessing, err := FindReadyPrompts(coll)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, false, err
	}

	if numProcessing > options.MaxStillProcessing {
		log.Info().Int("pending-prompts", len(prompts)).Int("max-processing", options.MaxStillProcessing).Int("num-processing", numProcessing).Msg(semLogContext + " max pending prompts exceeded")
		return nil, true, nil
	}

	if len(prompts) == 0 {
		return nil, false, nil
	}

	if len(prompts) > options.LimitCount && options.LimitCount > 0 {
		prompts = prompts[:options.LimitCount]
	}

	var totalWeight int
	var ndxLimit int
	for i := 0; i < len(prompts); i++ {
		totalWeight += int(prompts[i].Weight)
		if totalWeight > options.LimitWeight {
			break
		}
		ndxLimit = i + 1
	}

	if ndxLimit == 0 {
		if int(prompts[0].Weight) <= options.MaxLimitWeight || options.MaxLimitWeight <= 0 {
			prompts = prompts[:1]
		} else {
			err = errors.New("no prompts for weight limit reached")
			log.Error().Err(err).Msg(semLogContext + " no prompt to process")
			return nil, false, err
		}
	} else {
		prompts = prompts[:ndxLimit]
	}

	return prompts, false, nil
}

func FindReadyPrompts(coll *mongo.Collection) ([]PromptQueueItem, int, error) {
	const semLogContext = semLogPackageContext + "find-ready-prompts"

	filter := Filter{}
	filter.Or().AndEtEqTo(EntityType).AndStatusIn([]string{StatusReady, StatusProcessing})

	qr, err := FindByAggregationView(coll, &filter, true, new(options.FindOptions{}))
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return nil, -1, err
	}

	if len(qr.Data) == 0 {
		return nil, -1, nil
	}

	var numProcessing int
	var firstGroup PromptQueueItem
	var firstGroupFound bool
	for _, dto := range qr.Data {
		if dto.Status == StatusProcessing {
			numProcessing += int(dto.Count)
		} else {
			if !firstGroupFound {
				firstGroup = dto
				firstGroupFound = true
			}
		}
	}

	if firstGroupFound {
		filter = Filter{}
		filter.Or().AndEtEqTo(EntityType).AndStatusEqTo(StatusReady).
			AndDomainEqTo(firstGroup.Domain).AndSiteEqTo(firstGroup.Site).
			AndCategoryEqTo(firstGroup.Category)

		qr, err = Find(coll, &filter, false, options.Find())
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, -1, err
		}

		return qr.Data, numProcessing, nil
	}

	return nil, numProcessing, nil
}

func FindByAggregationView(collection *mongo.Collection, f *Filter, withCount bool, findOptions *options.FindOptions) (QueryResult, error) {
	const semLogContext = semLogPackageContext + "find-by-aggregation-view"

	fd := f.Build()
	log.Trace().Str("filter", util.MustToExtendedJsonString(fd, false, false)).Msg(semLogContext)

	qr := QueryResult{}

	// TODO
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	if withCount {
		countDocsOptions := options.Count()
		nr, err := collection.CountDocuments(ctx, fd, countDocsOptions)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return qr, err
		}

		qr.Records = int(nr)
		if nr == 0 {
			return qr, nil
		}
	}

	cur, err := cursorByAggregationView(collection, f, findOptions)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return qr, err
	}

	for cur.Next(context.Background()) {
		dto := PromptQueueItem{}
		err = cur.Decode(&dto)
		if err != nil {
			return qr, err
		}

		qr.Data = append(qr.Data, dto)
	}

	if cur.Err() != nil {
		return qr, cur.Err()
	}

	return qr, nil
}

func cursorByAggregationView(collection *mongo.Collection, f *Filter, findOptions *options.FindOptions) (*mongo.Cursor, error) {
	const semLogContext = semLogPackageContext + "cursor-by-aggregation-view"

	fd := f.Build()
	evtTraceLog := log.Trace().Str("filter", util.MustToExtendedJsonString(fd, false, false))
	evtErrLog := log.Error().Str("filter", util.MustToExtendedJsonString(fd, false, false))
	evtTraceLog.Msg(semLogContext)

	// TODO
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, bson.D{{"$match", fd}})
	if findOptions != nil {
		if findOptions.Skip != nil {
			pipeline = append(pipeline, bson.D{{"$skip", findOptions.Skip}})
		}
		if findOptions.Limit != nil {
			pipeline = append(pipeline, bson.D{{"$limit", findOptions.Limit}})
		}
	}
	pipeline = append(pipeline, bson.D{
		{"$group",
			bson.D{
				{"_id",
					bson.D{
						{"domain", "$domain"},
						{"site", "$site"},
						{"category", "$category"},
						{"status", "$status"},
					},
				},
				{"count", bson.D{{"$sum", 1}}},
				{"total_weight", bson.D{{"$sum", "$weight"}}},
			},
		},
	})

	pipeline = append(pipeline, bson.D{
		{"$project",
			bson.D{
				{"_id", 0},
				{"domain", "$_id.domain"},
				{"site", "$_id.site"},
				{"category", "$_id.category"},
				{"weight", "$total_weight"},
				{"count", "$count"},
				{"status", "$_id.status"},
			},
		},
	})

	opts := options.Aggregate()
	cur, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		evtErrLog.Err(err).Msg(semLogContext)
		return nil, err
	}

	for _, stage := range pipeline {
		b, err := bson.MarshalExtJSON(stage, true, true)
		if err == nil {
			fmt.Println(string(b))
		}
	}

	return cur, nil
}

// @tpm-schematics:end-region("bottom-file-section")
