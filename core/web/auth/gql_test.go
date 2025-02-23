package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/sessions/mocks"
	"github.com/smartcontractkit/chainlink/core/web/auth"
)

func Test_AuthenticateGQL_Unauthenticated(t *testing.T) {
	t.Parallel()

	sessionORM := &mocks.ORM{}
	sessionStore := sessions.NewCookieStore([]byte("secret"))

	r := gin.Default()
	r.Use(sessions.Sessions(auth.SessionName, sessionStore))
	r.Use(auth.AuthenticateGQL(sessionORM))

	r.GET("/", func(c *gin.Context) {
		user, ok := auth.GetGQLAuthenticatedUser(c)
		assert.False(t, ok)
		assert.Nil(t, user)

		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
}

func Test_AuthenticateGQL_Authenticated(t *testing.T) {
	t.Parallel()

	sessionORM := &mocks.ORM{}
	sessionStore := sessions.NewCookieStore([]byte(cltest.SessionSecret))
	sessionID := "sessionID"

	r := gin.Default()
	r.Use(sessions.Sessions(auth.SessionName, sessionStore))
	r.Use(auth.AuthenticateGQL(sessionORM))

	r.GET("/", func(c *gin.Context) {
		user, ok := auth.GetGQLAuthenticatedUser(c.Request.Context())
		assert.True(t, ok)
		assert.NotNil(t, user)

		c.String(http.StatusOK, "")
	})

	sessionORM.On("AuthorizedUserWithSession", sessionID).Return(clsessions.User{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	cookie := cltest.MustGenerateSessionCookie(t, sessionID)
	req.AddCookie(cookie)

	r.ServeHTTP(w, req)
}

func Test_GetAndSetGQLAuthenticatedUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := clsessions.User{}

	ctx = auth.SetGQLAuthenticatedUser(ctx, user)

	actual, ok := auth.GetGQLAuthenticatedUser(ctx)
	assert.True(t, ok)
	assert.Equal(t, &user, actual)
}
