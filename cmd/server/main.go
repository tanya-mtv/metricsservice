package main

<<<<<<< HEAD
import (
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/server"
)

func main() {
	// srv := server.NewServer(appLogger, cfg)
	cfg, err := config.InitConfig()
	if err != nil {
		// sugarLogger.Error("error initialazing config", zap.String("initConfig", "fail"), err)
		panic("error initialazing config")
	}
	appLogger := logger.NewAppLogger()
	srv := server.NewServer(appLogger, cfg)
	// appLogger.Fatal(srv.Run())
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
=======
func main() {}
>>>>>>> 707e40f (Initial commit)
