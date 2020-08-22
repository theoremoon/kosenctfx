package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

func (s *server) loginMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := s.getLoginUser(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": UnauthorizedMessage,
			})
		}
		return h(&loginContext{c, user})
	}
}

func (s *server) adminMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// tokenによる認証
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if strings.Contains(auth, s.Token) {
			user, err := s.app.GetAdminUser()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}

			return h(&loginContext{c, user})
		}

		// admin loginによる認証
		user, err := s.getLoginUser(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": UnauthorizedMessage,
			})
		}
		if !user.IsAdmin {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": AdminUnauthorizedMessage,
			})
		}
		return h(&loginContext{c, user})
	}
}

func (s *server) notLoginMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := s.getLoginUser(c)
		if user != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": AlreadyAuthorizedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfStartedMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if status == service.CTFNotStarted {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": CTFNotStartedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfNotStartedMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if status != service.CTFNotStarted {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": CTFAlreadyStartedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfRunningMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if status != service.CTFRunning {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
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
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": RegistrationClosedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfPlayableMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf, err := s.app.GetCTFConfig()
		if err != nil {
			return errorHandle(c, xerrors.Errorf(": %w", err))
		}

		if !conf.CTFOpen {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": CTFClosedMessage,
			})
		}
		return h(c)
	}
}
