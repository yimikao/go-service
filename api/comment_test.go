package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tutgo/business/comment"
	"tutgo/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateComment(t *testing.T) {
	tc := []struct {
		Name          string
		ReqBody       comment.NewComment
		StatusCode    int
		CallMuckFuncs func(*mocks.MockStorer)
	}{
		{
			Name:       "ok",
			ReqBody:    comment.NewComment{},
			StatusCode: http.StatusOK,
			CallMuckFuncs: func(cr *mocks.MockStorer) {
				cr.EXPECT().Create(gomock.Any()).Times(1).Return(&comment.Comment{}, nil)
			},
		},
	}

	for _, c := range tc {
		t.Run(c.Name, func(t *testing.T) {
			var (
				r *http.Request
				w http.ResponseWriter
			)

			b := new(bytes.Buffer)
			require.NoError(t, json.NewEncoder(b).Encode(c.ReqBody))
			r = httptest.NewRequest(http.MethodPost, "/test", b)
			w = httptest.NewRecorder()

			controller := gomock.NewController(t)
			defer controller.Finish()

			cmRepo := mocks.NewMockStorer(controller)

			ch := NewCommentHandler(cmRepo)
			c.CallMuckFuncs(cmRepo)

			ch.createComment(w, r)

			require.Equal(t, c.StatusCode, w.(*httptest.ResponseRecorder).Code)
		})
	}

}
