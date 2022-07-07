package session

import (
	"database/sql"
	"github.com/alexedwards/scs/postgresstore"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
	DBPool         *sql.DB
}

func (c *Session) InitSession() *scs.SessionManager {
	// how long should sessions last?
	minutes, err := strconv.Atoi(c.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	persist := strings.ToLower(c.CookiePersist) == "true"
	secure := strings.ToLower(c.CookieSecure) == "true"

	session := scs.New()

	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = c.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = c.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	// which session store?
	switch strings.ToLower(c.SessionType) {
	case "redis":
	case "mysql", "mariadb":
	case "postgres", "postgresql":
		session.Store = postgresstore.New(c.DBPool)
	default:
		//cookie
	}

	return session
}
