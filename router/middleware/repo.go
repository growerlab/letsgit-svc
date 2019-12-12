package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/growerlab/codev-svc/model"
	"github.com/pkg/errors"
)

func CtxRepoMiddleware(c *gin.Context) {
	if c.Request.URL.Path == "/graphql" {
		bodyRaw, err := c.GetRawData()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRaw))
		reqOptions := handler.NewRequestOptions(c.Request)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRaw))

		if strings.Contains(reqOptions.Query, "__type") {
			c.Next()
			return
		}

		reqRepo := getRepo(reqOptions)
		if reqRepo == nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("repo path is required"))
			return
		}

		repo, err := model.OpenRepo(reqRepo.Path, reqRepo.Name)
		if err == nil {
			c.Request.WithContext(context.WithValue(c, "repo", repo))
		} else {
			_ = c.AbortWithError(http.StatusNotFound, errors.Errorf("repo %s not found", reqRepo.fullPath()))
			return
		}
	}
	c.Next()
}

type repoRequest struct {
	Path string
	Name string
}

func (r *repoRequest) fullPath() string {
	return filepath.Join(r.Path, r.Name)
}

func getRepo(reqOptions *handler.RequestOptions) *repoRequest {
	var repoPath, repoName string

	repoPath, _ = reqOptions.Variables["path"].(string)
	repoName, _ = reqOptions.Variables["name"].(string)

	if len(repoPath) == 0 {
		return nil
	}

	return &repoRequest{
		Path: repoPath,
		Name: repoName,
	}
}
