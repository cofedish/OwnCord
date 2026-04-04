package admin_test

import (
	"net/http"
	"testing"

	"github.com/owncord/server/admin"
)

func TestAdminAPI_SessionExchange_SetsCookie(t *testing.T) {
	database := openAdminTestDB(t)
	token := createAdminUser(t, database)
	handler := admin.NewAdminAPI(database, "1.0.0", nil, nil, nil, nil)

	rr := doRequest(t, handler, http.MethodPost, "/session/exchange", token, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /session/exchange = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}

	found := false
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "owncord_admin_session" {
			found = true
			if cookie.HttpOnly != true {
				t.Error("admin session cookie should be HttpOnly")
			}
			if cookie.Path != "/admin" {
				t.Errorf("cookie path = %q, want /admin", cookie.Path)
			}
			if cookie.SameSite != http.SameSiteStrictMode {
				t.Errorf("cookie SameSite = %v, want Strict", cookie.SameSite)
			}
		}
	}
	if !found {
		t.Fatal("expected owncord_admin_session cookie to be set")
	}
}

func TestAdminAPI_Stats_AcceptsSessionCookie(t *testing.T) {
	database := openAdminTestDB(t)
	token := createAdminUser(t, database)
	handler := admin.NewAdminAPI(database, "1.0.0", nil, nil, nil, nil)

	req := newRequestWithCookie(t, http.MethodGet, "/stats", "owncord_admin_session", token, nil)
	w := serveAdminRequest(handler, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /stats with cookie = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestAdminAPI_SessionLogout_RevokesCookieSession(t *testing.T) {
	database := openAdminTestDB(t)
	token := createAdminUser(t, database)
	handler := admin.NewAdminAPI(database, "1.0.0", nil, nil, nil, nil)

	req := newRequestWithCookie(t, http.MethodPost, "/session/logout", "owncord_admin_session", token, nil)
	w := serveAdminRequest(handler, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("POST /session/logout = %d, want 204; body=%s", w.Code, w.Body.String())
	}

	req = newRequestWithCookie(t, http.MethodGet, "/stats", "owncord_admin_session", token, nil)
	w = serveAdminRequest(handler, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("GET /stats after logout = %d, want 401; body=%s", w.Code, w.Body.String())
	}
}
