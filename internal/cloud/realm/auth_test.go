package realm_test

import (
	"testing"

	"github.com/10gen/realm-cli/internal/cloud/realm"
	u "github.com/10gen/realm-cli/internal/utils/test"
	"github.com/10gen/realm-cli/internal/utils/test/assert"
	"github.com/10gen/realm-cli/internal/utils/test/mock"
)

func TestRealmAuthenticate(t *testing.T) {
	u.SkipUnlessRealmServerRunning(t)

	client := realm.NewClient(u.RealmServerURL())

	t.Run("Should fail with invalid credentials", func(t *testing.T) {
		_, err := client.Authenticate("username", "apiKey")
		assert.Equal(t,
			realm.ServerError{Message: "failed to authenticate with MongoDB Cloud API: You are not authorized for this resource."},
			err,
		)
	})

	t.Run("Should return session details with valid credentials", func(t *testing.T) {
		session, err := client.Authenticate(u.CloudUsername(), u.CloudAPIKey())
		assert.Nil(t, err)
		assert.NotEqual(t, "", session.AccessToken, "access token must not be blank")
		assert.NotEqual(t, "", session.RefreshToken, "refresh token must not be blank")
	})
}

func TestRealmAuthProfile(t *testing.T) {
	u.SkipUnlessRealmServerRunning(t)

	t.Run("Should fail without an auth client", func(t *testing.T) {
		client := realm.NewClient(u.RealmServerURL())

		_, err := client.AuthProfile()
		assert.Equal(t, realm.ErrInvalidSession{}, err)
	})

	t.Run("With an active session should return session details with valid credentials", func(t *testing.T) {
		client := newAuthClient(t)

		profile, err := client.AuthProfile()
		assert.Nil(t, err)
		assert.NotEqualf(t, 0, len(profile.Roles), "expected profile to have role(s)")
		assert.Equal(t, []string{u.CloudGroupID()}, profile.AllGroupIDs())
	})
}

func TestRealmAuthRefresh(t *testing.T) {
	u.SkipUnlessRealmServerRunning(t)

	t.Run("Does not refresh auth if request does not return invalid session code", func(t *testing.T) {
		client := realm.NewClient(u.RealmServerURL())

		session, err := client.Authenticate(u.CloudUsername(), u.CloudAPIKey())
		assert.Equal(t, nil, err)

		// invalidate the session's access token
		session.AccessToken = session.RefreshToken

		profile := mock.NewProfileWithSession(t, session)

		client = realm.NewAuthClient(profile.RealmBaseURL(), profile)
		_, err = client.AuthProfile()
		serverError, ok := err.(realm.ServerError)
		assert.True(t, ok, "expected %T to be server error", err)
		assert.Equal(t, realm.ServerError{Message: "invalid session: valid Issuer required"}, serverError)
	})

	t.Run("Should return the invalid session error when credentials are invalid", func(t *testing.T) {
		client := realm.NewClient(u.RealmServerURL())

		session, err := client.Authenticate(u.CloudUsername(), u.CloudAPIKey())
		assert.Equal(t, nil, err)

		// invalidate the session's tokens
		session.AccessToken = ""
		session.RefreshToken = session.AccessToken

		profile := mock.NewProfileWithSession(t, session)

		client = realm.NewAuthClient(profile.RealmBaseURL(), profile)
		_, err = client.AuthProfile()
		assert.Equal(t, realm.ErrInvalidSession{}, err)
	})
	// TODO: REALMC-7719 add test for expired credentials and test for ensuring profile cleared on invalid session
}

func newAuthClient(t *testing.T) realm.Client {
	t.Helper()

	client := realm.NewClient(u.RealmServerURL())

	session, err := client.Authenticate(u.CloudUsername(), u.CloudAPIKey())
	assert.Nil(t, err)

	profile := mock.NewProfileWithSession(t, session)

	return realm.NewAuthClient(profile.RealmBaseURL(), profile)
}