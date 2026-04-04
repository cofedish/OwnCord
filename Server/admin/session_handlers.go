package admin

import (
	"net/http"
	"strings"

	"github.com/owncord/server/auth"
	"github.com/owncord/server/db"
)

const adminSessionCookieName = "owncord_admin_session"

func isSecureRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	forwardedProto := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0])
	return strings.EqualFold(forwardedProto, "https")
}

func setAdminSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookieName,
		Value:    token,
		Path:     "/admin",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   isSecureRequest(r),
	})
}

func clearAdminSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookieName,
		Value:    "",
		Path:     "/admin",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   isSecureRequest(r),
		MaxAge:   -1,
	})
}

// handleExchangeSessionCookie upgrades a freshly issued Bearer token into an
// HttpOnly admin cookie so the browser UI does not need to persist tokens in
// localStorage. The raw bearer token must be supplied on this request because
// only the token hash is stored in the database.
func handleExchangeSessionCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := auth.ExtractBearerToken(r)
		if !ok {
			writeErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing or invalid authorization header")
			return
		}

		setAdminSessionCookie(w, r, token)
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

// handleSessionLogout clears the admin cookie and revokes the backing session.
func handleSessionLogout(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(adminSessionKey).(*db.Session)
		if !ok || sess == nil {
			clearAdminSessionCookie(w, r)
			writeErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
			return
		}

		if err := database.DeleteSession(sess.TokenHash); err != nil {
			writeErr(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to revoke session")
			return
		}

		clearAdminSessionCookie(w, r)
		w.WriteHeader(http.StatusNoContent)
	}
}
