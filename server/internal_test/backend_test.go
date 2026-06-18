package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ijry/sanjie/internal/app"
)

func perform(t *testing.T, application *app.App, method string, path string, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res := httptest.NewRecorder()
	application.Router.ServeHTTP(res, req)
	return res
}

func expectOK(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper()
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", res.Code, res.Body.String())
	}
}

func TestDashboardAndAlert(t *testing.T) {
	application := app.New(NewTestDB(t))
	res := perform(t, application, http.MethodGet, "/api/dashboard", "")
	expectOK(t, res)

	var body struct {
		Data struct {
			PendingCaptureCount int `json:"pendingCaptureCount"`
			UnhandledAlertCount int `json:"unhandledAlertCount"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Data.PendingCaptureCount != 1 || body.Data.UnhandledAlertCount != 3 {
		t.Fatalf("unexpected dashboard metrics: %+v", body.Data)
	}

	expectOK(t, perform(t, application, http.MethodPost, "/api/alerts/1/handle", ""))
}

func TestCaptureTaskFlow(t *testing.T) {
	application := app.New(NewTestDB(t))
	expectOK(t, perform(t, application, http.MethodPost, "/api/capture-tasks/1/start", ""))
	expectOK(t, perform(t, application, http.MethodPost, "/api/capture-tasks/1/complete", `{"note":"魂魄已签收"}`))

	res := perform(t, application, http.MethodGet, "/api/capture-tasks/1", "")
	expectOK(t, res)
	var body struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&body)
	if body.Data.Status != "completed" {
		t.Fatalf("status = %s", body.Data.Status)
	}
}

func TestLifeBookApprovalFlow(t *testing.T) {
	application := app.New(NewTestDB(t))
	expectOK(t, perform(t, application, http.MethodPost, "/api/life-book/1/freeze", ""))
	expectOK(t, perform(t, application, http.MethodPost, "/api/life-book/1/change-lifespan", `{"newLifespan":88,"reason":"救人一命"}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/approvals/1/approve", `{"note":"判官复核通过"}`))

	res := perform(t, application, http.MethodGet, "/api/life-book/1", "")
	expectOK(t, res)
	var body struct {
		Data struct {
			ActualLifespan int  `json:"actualLifespan"`
			Locked         bool `json:"locked"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&body)
	if body.Data.ActualLifespan != 88 || body.Data.Locked {
		t.Fatalf("unexpected life book state: %+v", body.Data)
	}
}

func TestReincarnationFlow(t *testing.T) {
	application := app.New(NewTestDB(t))
	expectOK(t, perform(t, application, http.MethodPost, "/api/reincarnations/1/approve", `{"note":"准了"}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/reincarnations/1/assign-quota", `{"quotaType":"普通胎"}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/reincarnations/1/soup", `{"inventoryId":1,"dose":1,"memoryAfter":0}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/reincarnations/1/complete", `{"note":"已入轮回"}`))

	res := perform(t, application, http.MethodGet, "/api/reincarnations/1", "")
	expectOK(t, res)
	var body struct {
		Data struct {
			Status     string `json:"status"`
			SoupStatus string `json:"soupStatus"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&body)
	if body.Data.Status != "completed" || body.Data.SoupStatus != "issued" {
		t.Fatalf("unexpected reincarnation state: %+v", body.Data)
	}
}

func TestHellWishSoupFlows(t *testing.T) {
	application := app.New(NewTestDB(t))
	expectOK(t, perform(t, application, http.MethodPost, "/api/hell/sentences/1/review", `{"note":"疑似判官手滑"}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/hell/floors/1/dispatch", `{"targetFloorId":2,"amount":40}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/wishes/1/route", `{"assignedDeity":"财神","note":"转财运组"}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/wishes/1/resolve", `{"note":"建议先上班"}`))
	expectOK(t, perform(t, application, http.MethodPost, "/api/soup/inventory/1/adjust", `{"delta":10,"note":"补货"}`))
	expectOK(t, perform(t, application, http.MethodGet, "/api/audit-logs", ""))
}
