package doc

import "go.uber.org/zap"

func initLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
