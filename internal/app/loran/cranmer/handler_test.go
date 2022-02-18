package cranmer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ali-a-a/loran/internal/app/loran/cranmer"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/ali-a-a/loran/pkg/router"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen,noctx
func TestAdd(t *testing.T) {
	t.Parallel()

	natsCfg := config.Default().NATS
	conn, err := cmq.CreateJetStreamConnection(natsCfg)
	assert.NoError(t, err)

	handler := cranmer.NewHandler(conn, natsCfg)

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

	e := router.New()

	e.POST("/api/add", handler.Add)

	go func() {
		err = e.Start(fmt.Sprintf(":%d", config.Default().Cranmer.Port))
		assert.NoError(t, err)
	}()

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
