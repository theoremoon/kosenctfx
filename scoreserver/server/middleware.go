package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theoremoon/kosenctfx/scoreserver/service"
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
				"messsage": AlreadyAuthorizedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfStartedMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, err)
		}

		if status == service.CTFNotStarted {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"messsage": CTFNotStartedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfNotStartedMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, err)
		}

		if status != service.CTFNotStarted {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"messsage": CTFAlreadyStartedMessage,
			})
		}
		return h(c)
	}
}

func (s *server) ctfRunningMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := s.app.CurrentCTFStatus()
		if err != nil {
			return errorHandle(c, err)
		}

		if status != service.CTFRunning {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"messsage": CTFNotRunningMessage,
			})
		}
		return h(c)
	}
}
