package services

import (
	"context"
	"fmt"

	"net/http"

	"reflect"

	"github.com/fpay/gopress"
	"github.com/sirupsen/logrus"
	"github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/models"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	// ElasticServiceName is the identity of elastic service
	ElasticServiceName = "elastic"
)

var fields logrus.Fields

// ElasticService type
type ElasticService struct {
	client *elastic.Client
	ctx    context.Context
	logger *gopress.Logger
}

// NewElasticService returns instance of elastic service
func NewElasticService(logger *gopress.Logger) *ElasticService {
	// TODO: 配置 elastic URL
	esURL := fmt.Sprintf("http://%s:%d", conf.GlobalConf.Elastic.Address, conf.GlobalConf.Elastic.Port)
	fields = logrus.Fields{
		"Service": "Elasticsearch",
		"URL":     esURL,
		"Error":   nil,
	}

	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetGzip(true),
		elastic.SetSniff(false),
		elastic.SetTraceLog(logger),
		elastic.SetInfoLog(logger),
		elastic.SetErrorLog(logger))
	if err != nil {
		logger.WithFields(fields).WithError(err).Error("Create elastic client failed")
		return nil
	} else if _, code, err := client.Ping(esURL).Do(ctx); err != nil || code != http.StatusOK {
		logger.WithFields(fields).WithField("code", code).WithError(err).Error("ping elastic failed")
		return nil
	}
	return &ElasticService{
		client: client,
		ctx:    ctx,
		logger: logger,
	}
}

// ServiceName is used to implements gopress.Service
func (s *ElasticService) ServiceName() string {
	return ElasticServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *ElasticService) RegisterContainer(c *gopress.Container) {
}

// Mapping action
func (s *ElasticService) Mapping(index string, mapping string) error {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := s.client.IndexExists(index).Do(s.ctx)
	if err != nil {
		s.logger.WithFields(fields).WithError(err).Error()
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := s.client.CreateIndex(index).BodyString(mapping).Do(s.ctx)
		if err != nil {
			s.logger.WithFields(fields).WithError(err).Error()
			return err
		}
		if !createIndex.Acknowledged {
			s.logger.WithFields(fields).Warnf("not acknowledged, more detail plz check es log")
		}
	}
	return nil
}

// Indexing action
func (s *ElasticService) Indexing(data interface{}, index, typ string) error {
	_, err := s.client.Index().
		Index(index).
		Type(typ).
		BodyJson(data).
		Do(s.ctx)
	if err != nil {
		// TODO: 缓存错误的文档索引
		s.logger.WithField("index", index).
			WithError(err).
			Errorf("doc: %#v")
	}
	return nil
}

// SearchPost action
// TODO: 考虑一下拓展搜索范围, 比如可以组合标题和文章主题
func (s *ElasticService) SearchPost(field, val string, test ...bool) ([]models.Post, error) {
	index := conf.PostIndex
	if len(test) > 0 {
		index = conf.TestIndex
	}
	matchQuery := elastic.NewMatchQuery(field, val)
	ret, err := s.client.Search().
		Index(index).
		Query(matchQuery).
		Sort("updated_at", true).
		From(0).Size(1000).
		Pretty(true).
		Do(s.ctx)
	if err != nil {
		s.logger.WithField("index", "post").WithError(err).Error()
		return nil, err
	}

	var results []models.Post
	for _, item := range ret.Each(reflect.TypeOf(models.Post{})) {
		if post, ok := item.(models.Post); ok {
			results = append(results, post)
		} else {
			s.logger.WithField("item", fmt.Sprintf("%#v", item)).Error("query post failed")
		}
	}
	return results, nil
}
