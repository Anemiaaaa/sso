package tests

import (
	ssov1 "github.com/Anemiaaaa/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sso/tests/suite"
	"testing"
	"time"
)

const (
	emptyAppId = 0
	appID      = 1
	appSecret  = "test_secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(
		ctx,
		&ssov1.RegisterRequest{
			Password: pass,
			Email:    email,
		},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})

	loginTime := time.Now()
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	uidRaw, ok := claims["uid"]
	require.True(t, ok, "uid claim is missing")

	emailRaw, ok := claims["email"]
	require.True(t, ok, "email claim is missing")

	appIDRaw, ok := claims["app_id"]
	require.True(t, ok, "app_id claim is missing")

	assert.Equal(t, respReg.GetUserId(), int64(uidRaw.(float64)))
	assert.Equal(t, email, emailRaw.(string))
	assert.Equal(t, appID, int(appIDRaw.(float64)))

	const deltaSecond = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSecond)

}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

const invalidAppID = 99999

func TestRegisterLogin_FailCases(t *testing.T) {
	t.Run("Register_EmptyEmail", func(t *testing.T) {
		ctx, st := suite.New(t)

		_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Password: randomFakePassword(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email")
	})

	t.Run("Register_EmptyPassword", func(t *testing.T) {
		ctx, st := suite.New(t)

		_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email: gofakeit.Email(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
	})

	t.Run("Register_EmailAlreadyExists", func(t *testing.T) {
		ctx, st := suite.New(t)

		email := gofakeit.Email()
		password := randomFakePassword()

		_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})
		require.NoError(t, err)

		_, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already")
	})

	t.Run("Login_NonExistentEmail", func(t *testing.T) {
		ctx, st := suite.New(t)

		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Email:    "notfound@example.com",
			Password: "pass1234",
			AppId:    appID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to login")
	})

	t.Run("Login_WrongPassword", func(t *testing.T) {
		ctx, st := suite.New(t)

		email := gofakeit.Email()
		pass := randomFakePassword()

		_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: pass,
		})
		require.NoError(t, err)

		_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Email:    email,
			Password: "wrong-password",
			AppId:    appID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to login")
	})

	t.Run("Login_EmptyEmail", func(t *testing.T) {
		ctx, st := suite.New(t)

		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Password: randomFakePassword(),
			AppId:    appID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email")
	})

	t.Run("Login_EmptyPassword", func(t *testing.T) {
		ctx, st := suite.New(t)

		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Email: gofakeit.Email(),
			AppId: appID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
	})

	t.Run("Login_InvalidAppID", func(t *testing.T) {
		ctx, st := suite.New(t)

		email := gofakeit.Email()
		pass := randomFakePassword()

		_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: pass,
		})
		require.NoError(t, err)

		_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Email:    email,
			Password: pass,
			AppId:    invalidAppID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to login")
	})
}
