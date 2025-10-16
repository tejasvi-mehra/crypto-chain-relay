package main

import (
	token "builder-router/internal/framework/jwk"
	"builder-router/internal/framework/logger"
	"builder-router/internal/framework/server"
	"builder-router/internal/service"
	"builder-router/internal/service/btc"
	"builder-router/internal/service/eth"
	"builder-router/internal/service/fiat"
	"builder-router/internal/service/mocks"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echo_log "github.com/labstack/gommon/log"
)

func main() {
	log := logger.GetLogger("main")
	config, err := server.NewServiceServerConfig()
	if err != nil {
		log.Panicw("failed to load config from env", "err", err)
	}
	// setup mock node
	mockNode := &mocks.Node{}
	// setup mock wallet
	wallet := &mocks.Wallet{}
	// tx-history
	txHistory := &mocks.TransactionHistory{}
	// should be separate service and should just be endpoint to service here
	ethBuilder := eth.NewEthBuilder(mockNode, wallet)
	btcBuilder := btc.NewBtcBuilder(mockNode, wallet)
	fiatBuilder := fiat.NewFiatBuilder()
	// router services
	ethRouter := eth.NewEthBuilderRouter(ethBuilder)
	btcRouter := btc.NewBtcBuilderRouter(btcBuilder)
	fiatRouter := fiat.NewFiatBuilderRouter(fiatBuilder)
	// service
	service := service.NewService(ethRouter, btcRouter, fiatRouter, txHistory)
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

func generateToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		accountID := c.QueryParam("accountId")
		if accountID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to make token, empty accountID"))
		}
		jwt, err := token.GenerateJwt("test", accountID, time.Duration(30*time.Minute))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to make token: %v", err))
		}
		return c.JSON(http.StatusOK, jwt)
	}
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
	e.GET("/token", generateToken())
	e.POST("/withdraw", func(c echo.Context) error {
		ctx := context.Background()
		var params service.InitiateWithdrawalParams
		if err := c.Bind(&params); err != nil {
			log.Errorw("cannot bind params", "err", err)
			return fmt.Errorf("failed to get params: %v", err)
		}
		res, err := s.InitiateWithdrawal(ctx, c, &params)
		if err != nil {
			return err
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return c.JSON(http.StatusOK, res)
	},
		// can add this for full echo as Middleware Pre check but since I'm using the same service to generate a token, added here
		token.Jwk(middleware.DefaultSkipper, nil),
	)
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
