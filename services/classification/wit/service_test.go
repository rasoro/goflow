package wit_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://api.wit.ai/message?v=20170307&q=book+flight+to+Quito": []httpx.MockResponse{
			httpx.NewMockResponse(200, `{"_text":"book flight to Quito","entities":{"intent":[{"confidence":0.84709152161066,"value":"book_flight"}]},"msg_id":"1M7fAcDWag76OmgDI"}`),
		},
	}))

	svc := wit.NewService(test.NewClassifier("Booking", "wit", []string{"book_flight", "book_hotel"}), "23532624376")

	eventLog := test.NewEventLog()

	classification, err := svc.Classify(session, "book flight to Quito", eventLog.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		flows.ExtractedIntent{Name: "book_flight", Confidence: decimal.RequireFromString(`0.84709152161066`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{}, classification.Entities)

	eventsJSON, _ := json.Marshal(eventLog.Events)
	test.AssertEqualJSON(t, []byte(`[
		{
			"classifier": {
				"name": "Booking",
				"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d"
			},
			"created_on": "2019-10-07T15:21:32.123456789Z",
			"http_logs": [
				{
					"created_on": "2019-10-07T15:21:30.123456789Z",
					"elapsed_ms": 1000,
					"request": "GET /message?v=20170307&q=book+flight+to+Quito HTTP/1.1\r\nHost: api.wit.ai\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer 23532624376\r\nAccept-Encoding: gzip\r\n\r\n",
					"response": "HTTP/1.0 200 OK\r\nContent-Length: 139\r\n\r\n{\"_text\":\"book flight to Quito\",\"entities\":{\"intent\":[{\"confidence\":0.84709152161066,\"value\":\"book_flight\"}]},\"msg_id\":\"1M7fAcDWag76OmgDI\"}",
					"status": "success",
					"url": "https://api.wit.ai/message?v=20170307&q=book+flight+to+Quito"
				}
			],
			"type": "classifier_called"
		}
	]`), eventsJSON, "events JSON mismatch")
}
