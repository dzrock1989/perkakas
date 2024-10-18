package syncservices

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"
	"github.com/tigapilarmandiri/perkakas/internal/models"
	"github.com/tigapilarmandiri/perkakas/internal/repositories/direktorat"
	"github.com/tigapilarmandiri/perkakas/internal/repositories/kepolisian"
	"github.com/tigapilarmandiri/perkakas/internal/repositories/pekerjaan"
	subdirektorat "github.com/tigapilarmandiri/perkakas/internal/repositories/sub_direktorat"
	"github.com/tigapilarmandiri/perkakas/internal/repositories/wilayah"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

func Begin(ctx context.Context, cl *kgo.Client, dbGorm *gorm.DB, ch chan<- struct{}) {
	// register topic
	adm := kadm.NewClient(cl)

	topics := strings.Split(configs.Config.Redpanda.Topics, ",")

	for _, v := range topics {
		_, err := adm.CreateTopics(ctx, 1, int16(configs.Config.Redpanda.ReplicaFactor), nil, v)
		if err != nil {
			util.Log.Error().Msg(err.Error())
		}
	}

	repoKepolisian := kepolisian.NewKepolisianRepository(dbGorm)
	repoDirektorat := direktorat.NewDirektoratRepository(dbGorm)
	repoSubDirektorat := subdirektorat.NewSubDirektoratRepository(dbGorm)
	repoWilayah := wilayah.NewWilayahRepository(dbGorm)
	repoPekerjaan := pekerjaan.NewPekerjaanRepository(dbGorm)

	var message string

	for {
		select {
		case <-ctx.Done():
			ch <- struct{}{}
			return
		default:
			fetches := cl.PollRecords(ctx, 1)
			if errs := fetches.Errors(); len(errs) > 0 {
				util.Log.Error().Msg(fmt.Sprint(errs))
				continue
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				if len(record.Headers) == 0 {
					util.Log.Error().Msg("len headers = 0")
					continue
				}
				action := record.Headers[0].Value
				message = fmt.Sprintf("topic: %s, action: %s, offset: %d", record.Topic, action, record.Offset)

			getTopic:
				switch record.Topic {
				case "kepolisian":
					var data models.Kepolisian

					switch string(action) {
					case "create":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoKepolisian.Create(ctx, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "update":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoKepolisian.Update(ctx, data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "delete":
						err := repoKepolisian.Delete(ctx, string(record.Value))
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "sync_wilayah":
						var datas []string
						err := json.Unmarshal(record.Value, &datas)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						if len(record.Headers) < 2 {
							util.Log.Error().Msg("len headers < 2, required uuid")
							break getTopic
						}
						bUuid := record.Headers[1].Value

						err = repoKepolisian.Tx(func(tx *kepolisian.Repository) error {
							err := tx.DeleteKepolisianIdOnTableJoin(ctx, string(bUuid))
							if err != nil {
								return err
							}
							return tx.SyncKepolisianHasWilayah(ctx, string(bUuid), datas)
						})
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}

					}
				case "direktorat":
					var data models.Direktorat

					switch string(action) {
					case "create":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoDirektorat.Create(ctx, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "update":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoDirektorat.Update(ctx, data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "delete":
						err := repoDirektorat.Delete(ctx, string(record.Value))
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}

					}
				case "sub_direktorat":
					var data models.SubDirektorat

					switch string(action) {
					case "create":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoSubDirektorat.Create(ctx, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "update":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoSubDirektorat.Update(ctx, data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "delete":
						err := repoSubDirektorat.Delete(ctx, string(record.Value))
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}

					}

				case "wilayah":
					var data models.Wilayah

					switch string(action) {
					case "create":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoWilayah.Create(ctx, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "update":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoWilayah.Update(ctx, data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "delete":
						err := repoWilayah.Delete(ctx, string(record.Value))
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					}

				case "pekerjaan":
					var data models.Pekerjaan

					switch string(action) {
					case "create":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoPekerjaan.Create(ctx, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "update":
						err := json.Unmarshal(record.Value, &data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
						err = repoPekerjaan.Update(ctx, data)
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					case "delete":
						err := repoPekerjaan.Delete(ctx, string(record.Value))
						if err != nil {
							util.Log.Err(err).Msg(message)
							break getTopic
						}
					}
				}
			}
			err := cl.CommitRecords(ctx, fetches.Records()...)
			if err != nil {
				util.Log.Err(err).Msg("failed to commit to broker")
			}

		}
	}
}
