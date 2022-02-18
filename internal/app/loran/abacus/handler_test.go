package abacus_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/internal/app/loran/abacus/model"
	"github.com/ali-a-a/loran/internal/app/loran/cranmer"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/ali-a-a/loran/pkg/redis"
	"github.com/stretchr/testify/assert"

	"github.com/ali-a-a/loran/internal/app/loran/abacus"
)

// Note: Since I want to test the whole abacus package, I don't mock redis or nats clients.
//nolint:funlen
func TestRun(t *testing.T) {
	t.Parallel()

	natsCfg := config.Default().NATS
	conn, err := cmq.CreateJetStreamConnection(natsCfg)
	assert.NoError(t, err)

	redisCfg := config.Default().Redis
	rc, err := redis.NewRedisClient(redisCfg)
	assert.NoError(t, err)

	cr := model.NewInMemoryCalculator(rc)

	handler, err := abacus.NewHandler(cr, conn, natsCfg)
	assert.NoError(t, err)

	cases := []struct {
		name        string
		messages    []cranmer.AddRequest
		target      cranmer.AddRequest
		expectedCnt int64
	}{
		{
			name: "one publish same as target",
			messages: []cranmer.AddRequest{
				{
					UserID:   123,
					EntityID: 321,
				},
			},
			target: cranmer.AddRequest{
				UserID:   123,
				EntityID: 321,
			},
			expectedCnt: 1,
		},
		{
			name: "two publish two same as target",
			messages: []cranmer.AddRequest{
				{
					UserID:   567,
					EntityID: 123,
				},
				{
					UserID:   567,
					EntityID: 123,
				},
			},
			target: cranmer.AddRequest{
				UserID:   567,
				EntityID: 123,
			},
			expectedCnt: 1,
		},
		{
			name: "two publish one same as target",
			messages: []cranmer.AddRequest{
				{
					UserID:   569,
					EntityID: 789,
				},
				{
					UserID:   560,
					EntityID: 789,
				},
			},
			target: cranmer.AddRequest{
				UserID:   560,
				EntityID: 789,
			},
			expectedCnt: 2,
		},
	}

	for i := range cases {
		test := cases[i]

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			go func() {
				err = handler.Run()
				assert.NoError(t, err)
			}()

			time.Sleep(1 * time.Second)

			for i := range test.messages {
				b, err := json.Marshal(&test.messages[i])
				assert.NoError(t, err)

				_, err = conn.JS.Publish(natsCfg.JetStream.Consumer.Subject, b)
				assert.NoError(t, err)
			}

			cnt, err := cr.Count(context.Background(), test.target.EntityID)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCnt, cnt)
		})
	}
}
