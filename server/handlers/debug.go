package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx"
)

type debugHandler struct {
	http.Handler
}

func DebugHandler() debugHandler {
	return debugHandler{}
}

type ServiceGroup struct {
	ID                  int64  `json:"id"`
	GroupName           string `json:"group_name"`
	TemplateID          int64  `json:"template_id"`
	AccountID           string `json:"account_id"`
	Capacity            int    `json:"capacity"`
	HealthCheckInterval int    `json:"health_check_interval"`
}

func FindGroupByID(ctx context.Context, key int64, accountId string) (*ServiceGroup, error) {
	var group ServiceGroup

	sqlStatement := `SELECT id, name, account_id, template_id, capacity, health_check_interval
FROM triton.tsg_groups
WHERE account_id = $2 and id = $1
AND archived = false;`

	db, ok := GetDBPool(ctx)
	if !ok {
		return nil, ErrNoConnPool
	}

	err := db.QueryRowEx(ctx, sqlStatement, nil, key, accountId).
		Scan(&group.ID,
			&group.GroupName,
			&group.AccountID,
			&group.TemplateID,
			&group.Capacity,
			&group.HealthCheckInterval)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (h debugHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	ctx := req.Context()

	grp, err := FindGroupByID(ctx, int64(320377354180919298), "joyent")
	if err == pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = fmt.Fprintf(w, "It's happening to %s!\n", grp.GroupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session, ok := GetAuthSession(ctx)
	if !ok {
		// NOTE(justinwr): This should NEVER happen since it is caught by our
		// upstream AuthHandler.
		http.Error(w, ErrNoSession.Error(), http.StatusUnauthorized)
	}

	_, err = fmt.Fprintf(w, "Session AccountID == %q!\n", session.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
