package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/loader"
	"github.com/theoremoon/kosenctfx/scoreserver/resolver"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
	"golang.org/x/xerrors"
)

/// Attach Logged in Team or nil into Context for GraphQL endpoint
func (s *server) resolveLoginMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// tokenによる認証
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if strings.Contains(auth, s.Token) {
			t, err := s.app.GetAdminTeam()
			if err != nil {
				return errorHandle(c, xerrors.Errorf(": %w", err))
			}
			r := c.Request()
			c.SetRequest(r.WithContext(resolver.AttachTeam(r.Context(), t)))
			return h(c)
		}

		// ログインによる認証
		t, _ := s.getLoginTeam(c)
		r := c.Request()
		c.SetRequest(r.WithContext(resolver.AttachTeam(r.Context(), t)))
		return h(c)
	}
}

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
		return h(&loginContext{c, team})
	}
}

func (s *server) adminMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// tokenによる認証
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if strings.Contains(auth, s.Token) {
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
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
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
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
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
