package doc

import (
	"github.com/holive/doc/app/config"
	"github.com/holive/doc/app/mongo"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Doc struct {
	Cfg      *config.Config
	Services *Services
}

type Services struct {
}

func New() (*Doc, error) {
	var (
		err error
		f   = &Doc{}
	)

	f.Cfg, err = loadConfig("./config")
	if err != nil {
		return nil, errors.Wrap(err, "could not load config")
	}

	db, err := initMongoClient(f.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize mongo client")
	}

	httpClient, err := initHTTPClient(f.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize http client")
	}

	logger, err := initLogger()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize logger")
	}

	f.Services = initServices(f.Cfg, db, httpClient, logger)

	return f, nil
}

func initServices(cfg *config.Config, db *mongo.Client, client infraHTTP.Runner, logger *zap.SugaredLogger) *Services {


	return &Services{

	}
}
