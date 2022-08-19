package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/gomodule/redigo/redis"

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
	RedisPool      *redis.Pool
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
		session.Store = redisstore.New(c.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(c.DBPool)
	case "postgres", "postgresql":
		session.Store = postgresstore.New(c.DBPool)
	default:
		//cookie
	}

	return session
}
