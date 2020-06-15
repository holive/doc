package doc

import (
	"github.com/holive/doc/app/config"
	"github.com/holive/doc/app/docApi"
	"github.com/holive/doc/app/mongo"
	"github.com/holive/doc/app/squads"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Doc struct {
	Cfg      *config.Config
	Services *Services
}

type Services struct {
	DocApi *docApi.Service
	Squads *squads.Service
}

func New() (*Doc, error) {
	var (
		err error
		f   = &Doc{}
	)

	f.Cfg, err = loadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not load config")
	}

	db, err := initMongoClient(f.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize mongo client")
	}

	logger, err := initLogger()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize logger")
	}

	f.Services = initServices(db, logger)

	return f, nil
}

func initServices(db *mongo.Client, logger *zap.SugaredLogger) *Services {
	docApiService := initDocApiService(db)
	squadsService := initSquadsService(db)

	return &Services{
		DocApi: docApiService,
		Squads: squadsService,
	}
}
