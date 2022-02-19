package cranmer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/internal/app/loran/cranmer"
	"github.com/ali-a-a/loran/internal/app/loran/model"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/ali-a-a/loran/pkg/redis"
	"github.com/ali-a-a/loran/pkg/router"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	natsCfg := config.Default().NATS
	conn, _ := cmq.CreateJetStreamConnection(natsCfg)

	redisCfg := config.Default().Redis
	rc, _ := redis.NewRedisClient(redisCfg)

	cr := model.NewInMemoryCalculator(rc)

	_ = cr.Add(context.Background(), 1, 234)
	_ = cr.Add(context.Background(), 2, 234)

	handler := cranmer.NewHandler(conn, cr, natsCfg)
	e := router.New()

	e.POST("/api/add", handler.Add)
	e.POST("/api/count", handler.Count)

	go func() {
		_ = e.Start(fmt.Sprintf(":%d", config.Default().Cranmer.Port))
	}()

	os.Exit(m.Run())
}

//nolint:funlen,noctx
func TestAdd(t *testing.T) {
	t.Parallel()

	req := &cranmer.AddRequest{
		UserID:   123,
		EntityID: 234,
	}

	invalidReq := map[string]int{
		"x": 123,
		"y": 234,
	}

	cases := []struct {
		name string
		req  interface{}
		code int
		fail bool
	}{
		{
			name: "successful",
			req:  req,
			code: http.StatusOK,
			fail: false,
		},
		{
			name: "invalid request",
			req:  invalidReq,
			code: http.StatusBadRequest,
			fail: true,
		},
	}

	for i := range cases {
		test := cases[i]

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(test.req)
			assert.NoError(t, err)

			resp, err := http.Post(fmt.Sprintf("http://localhost:%d/api/add", config.Default().Cranmer.Port),
				"application/json", bytes.NewReader(b))

			defer func() {
				err = resp.Body.Close()
				assert.NoError(t, err)
			}()

			assert.NoError(t, err)
			if test.fail {
				assert.Equal(t, test.code, resp.StatusCode)
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, test.code, resp.StatusCode)
				assert.Equal(t, "{\"status\":\"ok\"}\n", string(body))
			}
		})
	}
}

//nolint:funlen,noctx
func TestCount(t *testing.T) {
	t.Parallel()

	req := &cranmer.CountRequest{
		EntityID: 234,
	}

	invalidReq := map[string]int{
		"x": 123,
	}

	cases := []struct {
		name        string
		req         interface{}
		code        int
		expectedCnt int
		fail        bool
	}{
		{
			name:        "successful",
			req:         req,
			code:        http.StatusOK,
			expectedCnt: 2,
			fail:        false,
		},
		{
			name: "invalid request",
			req:  invalidReq,
			code: http.StatusBadRequest,
			fail: true,
		},
	}

	for i := range cases {
		test := cases[i]

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(test.req)
			assert.NoError(t, err)

			resp, err := http.Post(fmt.Sprintf("http://localhost:%d/api/count", config.Default().Cranmer.Port),
				"application/json", bytes.NewReader(b))

			defer func() {
				err = resp.Body.Close()
				assert.NoError(t, err)
			}()

			assert.NoError(t, err)
			if test.fail {
				assert.Equal(t, test.code, resp.StatusCode)
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, test.code, resp.StatusCode)
				assert.Equal(t, fmt.Sprintf("{\"count\":%d,\"status\":\"ok\"}\n",
					test.expectedCnt), string(body))
			}
		})
	}
}
