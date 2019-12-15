package oauth

import "testing"

func TestSimpleAuthenticationProvider_Init(t *testing.T) {

	auth := SimpleAuthenticationProvider{}

	// since no file exists ok should be false
	if ok, err := auth.Init("test.json"); err != nil || ok != false {
		t.Error(err)
	}

	if err := auth.RegisterClient("123", "456"); err != nil {
		t.Error(err)
	}

	if err := auth.RegisterAgent("test", "password"); err != nil {
		t.Error(err)
	}

	auth2 := SimpleAuthenticationProvider{}
	if ok, err := auth2.Init("test.json"); err != nil || ok != true {
		t.Error(err)
	}

	if auth2.db.Clients[0].ID != auth2.db.Clients[0].ID {
		t.Error("client id does not match")
	}

	if auth2.db.Clients[0].Secret != auth2.db.Clients[0].Secret {
		t.Error("client secret does not match")
	}

	if auth2.db.Agents[0].ID != auth2.db.Agents[0].ID {
		t.Error("agent id does not match")
	}

	if auth2.db.Agents[0].Password != auth2.db.Agents[0].Password {
		t.Error("agent password does not match")
	}

	if code, err := auth2.GenerateAuthCode("123", "test"); err == nil {
		if token, err := auth2.GenerateToken("123", "456", code); err == nil {
			if agent, refresh, err := auth2.GetAgentUserIDForToken(token.AccessToken); err == nil {
				if refresh {
					t.Error("expected access token, got refresh token")
				}
				if agent != "test" {
					t.Error("expected test got: " + agent)
				}
			} else {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}

	if err := auth2.Delete(); err != nil {
		t.Error(err)
	}

}
