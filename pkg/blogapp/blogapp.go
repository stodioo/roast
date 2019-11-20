package blogapp

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/stodioo/roast/pkg/blogcore"
)

func NewBlogApp(envVarPrefix string) (*blogcore.BlogCore, error) {
	var err error
	var blogDB *gorm.DB
	cfg := ParseConfigFromEnv(envVarPrefix)

	if cfg.BlogDBUrl == "" {
		return nil, errors.New("envvar " + envVarPrefix + "_DB_URL not provided")
	}

	blogDB, err = connectToPostgres(cfg.BlogDBUrl, 4, 16)

	if err != nil {
		return nil, err
	}
	blogCore := blogcore.NewBlogCore(blogDB)

	return blogCore, nil
}

func connectToPostgres(dbURL string, defaultMaxIdleConns, defaultMaxOpenConns int) (*gorm.DB, error) {
	var db *gorm.DB
	parsedDBURL, err := url.Parse(dbURL)

	if err != nil {
		return nil, err
	}

	maxIdleConns := int64(defaultMaxIdleConns)
	maxOpenConns := int64(defaultMaxOpenConns)
	queryPart := parsedDBURL.Query()
	if maxIdleConnsStr := queryPart.Get("max_idle_conns"); maxIdleConnsStr != "" {
		queryPart.Del("max_idle_conns")
		maxIdleConns, err = strconv.ParseInt(maxIdleConnsStr, 10, 32)
		if err != nil {
			return nil, errors.New("Unable to parse max_idle_conns from query parameter")
		}
	}

	if maxOpenConnsStr := queryPart.Get("max_open_conns"); maxOpenConnsStr != "" {
		queryPart.Del("max_open_conns")
		maxOpenConns, err = strconv.ParseInt(maxOpenConnsStr, 10, 32)
		if err != nil {
			return nil, errors.New("Unable to parse max_open_conns from query parameters")
		}
	}

	if maxIdleConns == 0 {
		maxIdleConns = 4
	}

	if maxOpenConns == 0 {
		maxOpenConns = 16
	}

	parsedDBURL.RawQuery = queryPart.Encode()
	dbURL = parsedDBURL.String()

	for {
		db, err = gorm.Open("postgres", dbURL)

		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "connect: connection timeout") {
			return nil, err
		}

		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)
	}

	if db != nil {
		db.DB().SetMaxIdleConns(int(maxIdleConns))
		db.DB().SetMaxOpenConns(int(maxOpenConns))
	}

	return db, nil
}
