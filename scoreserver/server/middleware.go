package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/loader"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

func (s *server) attachLoaderMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		c.SetRequest(r.WithContext(loader.Attach(s.app, r.Context())))
		return h(c)
	}
}

func (s *server) loginMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		team, err := s.getLoginTeam(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": UnauthorizedMessage,
			})
		}

		// ログイン通っているときアクティブなトークンの数を記録する
		token, err := s.getLoginToken(c)
		if err == nil {
			// tokenそのままを使うのは嫌だけどだいたい単射であってほしい
			tokenHash := sha256.Sum256([]byte(token))
			tokenStr := hex.EncodeToString(tokenHash[:])

			// 値は上書きしたいのでmemberの値で一度削除する
			s.redis.ZRem(context.Background(), sessionSetKey, tokenStr)
			s.redis.ZAdd(context.Background(), sessionSetKey, &redis.Z{
				Score:  float64(time.Now().Add(sessionActiveDuration).Unix()),
				Member: tokenStr,
			})
		}
		return h(&loginContext{c, team})
	}
}

func (s *server) adminMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorized := false
		// tokenによる認証
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if strings.HasPrefix(auth, "Bearer ") {
			if auth[len("Bearer "):] == s.Token {
				authorized = true
			}
		}

		// basic認証
		_, password, isBasic := c.Request().BasicAuth()
		if isBasic {
			if password == s.Token {
				authorized = true
			}
		}

		if authorized {
			team, err := s.app.GetAdminTeam()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			return h(&loginContext{c, team})
		}

		// admin loginによる認証
		team, err := s.getLoginTeam(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": UnauthorizedMessage,
			})
		}
		if !team.IsAdmin {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": AdminUnauthorizedMessage,
			})
		}
		return h(&loginContext{c, team})
	}
}

func (s *server) notLoginMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		team, _ := s.getLoginTeam(c)
		if team != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": AlreadyAuthorizedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfStartedMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		status := service.CalcCTFStatus(conf)

		if status == service.CTFNotStarted {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": CTFNotStartedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfNotStartedMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		status := service.CalcCTFStatus(conf)

		if status != service.CTFNotStarted {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": CTFAlreadyStartedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfRunningMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}
		status := service.CalcCTFStatus(conf)

		if status != service.CTFRunning {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": CTFNotRunningMessage,
			})
		}
		return h(c)
	}
}

func (s *server) registerableMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if !conf.RegisterOpen {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": RegistrationClosedMessage,
			})
		}
		return h(c)
	}
}
