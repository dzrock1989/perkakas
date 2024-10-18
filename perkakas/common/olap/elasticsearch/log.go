package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/tigapilarmandiri/perkakas/configs"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type historyData struct {
	Env           string    `json:"env"`
	ExecutedAt    time.Time `json:"executed_at"`
	PIC           string    `json:"pic"`
	ServiceName   string    `json:"service_name"`
	TableName     string    `json:"table_name" gorm:"column:table_name"`
	EventId       string    `json:"event_id"`
	Operation     string    `json:"operation"`
	Payload       any       `json:"payload" gorm:"-"`
	PayloadString string    `json:"-" gorm:"column:payload"`
}

func CreateHistory(executedAt time.Time, pic, serviceName, tableName, eventId string, payload any) historyData {
	return historyData{
		EventId:     eventId,
		Env:         configs.Config.Env,
		ExecutedAt:  time.Now(),
		PIC:         pic,
		ServiceName: strings.ToLower(serviceName),
		TableName:   strings.ToLower(tableName),
		Operation:   "create",
		Payload:     payload,
	}
}

func UpdateHistory(executedAt time.Time, pic, serviceName, tableName, eventId string, payload any) historyData {
	return historyData{
		EventId:     eventId,
		Env:         configs.Config.Env,
		ExecutedAt:  time.Now(),
		PIC:         pic,
		ServiceName: strings.ToLower(serviceName),
		TableName:   strings.ToLower(tableName),
		Operation:   "update",
		Payload:     payload,
	}
}

func DeleteHistory(executedAt time.Time, pic, serviceName, tableName, eventId string, payload any) historyData {
	return historyData{
		EventId:     eventId,
		Env:         configs.Config.Env,
		ExecutedAt:  time.Now(),
		PIC:         pic,
		ServiceName: strings.ToLower(serviceName),
		TableName:   strings.ToLower(tableName),
		Operation:   "delete",
		Payload:     payload,
	}
}

var onceHistory sync.Once

var client *elasticsearch.Client

func StoreHistory(ctx context.Context, db *gorm.DB, index string, data historyData) (err error) {
	defer func() {
		if err != nil {
			apm.CaptureError(ctx, err)
		}
	}()

	if db != nil {
		b, err := json.Marshal(data.Payload)
		if err != nil {
			apm.CaptureError(ctx, err)
		}
		data.PayloadString = string(b)
		return db.Table("audit_logs").Create(&data).Error
	}

	onceHistory.Do(func() {
		client, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses:  []string{configs.Config.Elastic.Url},
			APIKey:     configs.Config.Elastic.ApiKey,
			MaxRetries: 3,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		})
	})
	if err != nil {
		return err
	}

	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r, err := client.Index(index, bytes.NewReader(bData))
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return nil
}
