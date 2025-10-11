package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wallet/internal/framework/db"
	"wallet/internal/framework/logger"
	"wallet/internal/framework/server"
	"wallet/internal/service"
	"wallet/internal/service/kms"
	"wallet/internal/service/repo"

	"github.com/labstack/echo/v4/middleware"
	echo_log "github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
)

func main() {
	log := logger.GetLogger("main")
	config, err := server.NewServiceServerConfig()
	if err != nil {
		log.Panicw("failed to load config from env", "err", err)
	}
	// setup db
	dbConfig, err := db.NewConfig()
	if err != nil {
		log.Panicw("failed to get db config", "err", err)
	}
	db, err := db.NewDb(dbConfig)
	if err != nil {
		log.Panicw("failed to setup db", "err", err)
	}
	repo := repo.NewRepo(db)
	// setup kms
	kmsConfig, err := kms.NewKMSConfig()
	if err != nil {
		log.Panicw("failed to load KMS config", "err", err)
	}
	kms, err := kms.NewKMS(kmsConfig)
	if err != nil {
		log.Panicw("failed to load KMS", "err", err)
	}
	// setup service
	service := service.NewService(kms, repo)
	e, err := ServiceEcho(service)
	if err != nil {
		log.Panicw("error", "err", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := server.Serve(e, config.ToServerConfig())
	<-quit
	log.Infow("shutting down")
	s.Shutdown()
}

func ServiceEcho(s *service.Service) (*echo.Echo, error) {
	log := logger.GetLogger("ServiceEcho")

	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  echo_log.ERROR,
	}))
	e.Pre(middleware.BodyLimit("1000K"))
	// ideally here we should also setup some JWT verification

	// setup endpoints
	e.POST("/address", func(c echo.Context) error {
		ctx := context.Background()
		var params service.GetAddressPublicKeyParams
		if err := c.Bind(&params); err != nil {
			log.Errorw("cannot bind params", "err", err)
			return fmt.Errorf("failed to get params: %v", err)
		}
		res, err := s.GetAddressPublicKey(ctx, &params)
		if err != nil {
			return err
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return c.JSON(http.StatusOK, res)
	})
	e.POST("/buildTxEth", func(c echo.Context) error {
		ctx := context.Background()
		var params service.BuildTxEthParams
		if err := c.Bind(&params); err != nil {
			log.Errorw("cannot bind params", "err", err)
			return fmt.Errorf("failed to get params: %v", err)
		}
		res, err := s.BuildTxEth(ctx, &params)
		if err != nil {
			return err
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return c.JSON(http.StatusOK, res)
	})
	e.POST("/signTx", func(c echo.Context) error {
		ctx := context.Background()
		var params service.SignTxParams
		if err := c.Bind(&params); err != nil {
			log.Errorw("cannot bind params", "err", err)
			return fmt.Errorf("failed to get params: %v", err)
		}
		res, err := s.SignTx(ctx, &params)
		if err != nil {
			return err
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return c.JSON(http.StatusOK, res)
	})
	e.HTTPErrorHandler = ErrorHandler(e.DefaultHTTPErrorHandler)

	return e, nil
}

func ErrorHandler(next echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// same check as default handler
		if c.Response().Committed {
			return
		}

		type Err struct {
			Message string `json:"message"`
		}

		if err != nil {
			errRes := Err{
				Message: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, errRes)
			return
		}

		if next != nil {
			next(err, c)
		}
	}
}
