package middleware

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// SentinelConfig defines the config for Sentinel middleware.
	SentinelConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// AllowOrigin defines a list of origins that may access the resource.
		// Optional. Default value []string{"*"}.
		AllowOrigins []string `json:"allow_origins"`
		// AllowsAllOrigin defines all origin should be accept or not
		// Optional. Default false, but if AllowOrigins contains "*", automatically turn this true.
		AllowsAllOrigin      bool             `json:"allows_all_origin"`
		AllowsStrictOrigins  []string         `json:"-"`
		AllowsStrictMatchers []*regexp.Regexp `json:"-"`
	}
)

var (
	// DefaultSentinelConfig is the default Sentinel middleware config.
	DefaultSentinelConfig = SentinelConfig{
		Skipper:         middleware.DefaultSkipper,
		AllowOrigins:    []string{"*"},
		AllowsAllOrigin: true,
	}
)

// Sentinel returns a Cross-Origin Resource Sharing (Sentinel) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_Sentinel
func Sentinel() echo.MiddlewareFunc {
	return SentinelWithConfig(DefaultSentinelConfig)
}

// SentinelWithConfig returns a Sentinel middleware with config.
// See: `Sentinel()`.
func SentinelWithConfig(config SentinelConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultSentinelConfig.Skipper
	}
	allowedOrigins := make([]string, 0)
	allowedMatchers := make([]*regexp.Regexp, 0)
	if len(config.AllowOrigins) != 0 {
		for _, v := range config.AllowOrigins {
			switch {
			case v == "*":
				config.AllowsAllOrigin = true
			case strings.ContainsAny(v, "^[]()?*+$"):
				allowedMatchers = append(allowedMatchers, regexp.MustCompile(v))
			default:
				allowedOrigins = append(allowedOrigins, v)
			}
		}
	} else {
		config.AllowOrigins = DefaultSentinelConfig.AllowOrigins
	}
	config.AllowsStrictOrigins = allowedOrigins
	config.AllowsStrictMatchers = allowedMatchers

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) || config.AllowsAllOrigin {
				return next(c)
			}

			req := c.Request()
			origin := req.Header.Get(echo.HeaderOrigin)

			// Check allowed origins
			for _, o := range allowedOrigins {
				if o == origin {
					return next(c)
				}
			}

			// Check wildcard domains
			for i := range allowedMatchers {
				o := allowedMatchers[i]
				if o.MatchString(origin) {
					return next(c)
				}
			}
			log.Printf("Accessed Origin: %s is not allowed to enter the hatch.\n", origin)
			return c.NoContent(http.StatusNoContent)
		}
	}
}
