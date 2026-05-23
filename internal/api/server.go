package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"atp-services/internal/app"
	"atp-services/internal/models"
)

type Server struct {
	app    *app.Application
	static string
}

func NewServer(application *app.Application, staticDir string) *Server {
	return &Server{app: application, static: staticDir}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/login", s.handleLogin)
	mux.HandleFunc("/api/logout", s.handleLogout)
	mux.HandleFunc("/api/me", s.auth(s.handleMe))
	mux.HandleFunc("/api/clients", s.auth(s.handleClients))
	mux.HandleFunc("/api/vehicles", s.auth(s.handleVehicles))
	mux.HandleFunc("/api/tariffs", s.auth(s.handleTariffs))
	mux.HandleFunc("/api/users", s.auth(s.handleUsers))
	mux.HandleFunc("/api/orders", s.auth(s.handleOrders))
	mux.HandleFunc("/api/orders/status", s.auth(s.handleOrderStatus))
	mux.HandleFunc("/api/orders/preview-price", s.auth(s.handlePreviewPrice))
	mux.HandleFunc("/api/schedule", s.auth(s.handleSchedule))
	mux.HandleFunc("/api/waybills", s.auth(s.handleWaybills))
	mux.HandleFunc("/api/shift/open", s.auth(s.handleOpenShift))
	mux.HandleFunc("/api/shift/close", s.auth(s.handleCloseShift))
	mux.HandleFunc("/api/shift/status", s.auth(s.handleShiftStatus))
	mux.HandleFunc("/api/drivers/available", s.auth(s.handleDriversAvailable))
	mux.HandleFunc("/api/dashboard", s.auth(s.handleDashboard))
	mux.HandleFunc("/api/reports/drivers", s.auth(s.handleDriverRating))
	mux.HandleFunc("/api/audit", s.auth(s.handleAudit))

	if s.static != "" {
		fs := http.FileServer(http.Dir(s.static))
		mux.Handle("/", spaHandler(s.static, fs))
	}
	return cors(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.LoginRequest
	if !readJSON(r, &req) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	resp, err := s.app.Login(req)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, resp)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	_ = s.app.Logout(token)
	writeJSON(w, map[string]bool{"ok": true})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	u, err := s.app.Me(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, u)
}

func (s *Server) handleClients(w http.ResponseWriter, r *http.Request) {
	token := tokenFrom(r)
	switch r.Method {
	case http.MethodGet:
		list, err := s.app.ListClients(token)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, list)
	case http.MethodPost:
		var c models.Client
		if !readJSON(r, &c) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		out, err := s.app.SaveClient(token, c)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, out)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleVehicles(w http.ResponseWriter, r *http.Request) {
	token := tokenFrom(r)
	switch r.Method {
	case http.MethodGet:
		list, err := s.app.ListVehicles(token)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, list)
	case http.MethodPost:
		var v models.Vehicle
		if !readJSON(r, &v) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		out, err := s.app.SaveVehicle(token, v)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, out)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTariffs(w http.ResponseWriter, r *http.Request) {
	token := tokenFrom(r)
	switch r.Method {
	case http.MethodGet:
		list, err := s.app.ListTariffs(token)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, list)
	case http.MethodPost:
		var t models.Tariff
		if !readJSON(r, &t) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		out, err := s.app.SaveTariff(token, t)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, out)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	token := tokenFrom(r)
	switch r.Method {
	case http.MethodGet:
		list, err := s.app.ListUsers(token)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, list)
	case http.MethodPost:
		var body struct {
			User     models.User `json:"user"`
			Password string      `json:"password"`
		}
		if !readJSON(r, &body) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		out, err := s.app.CreateUser(token, body.User, body.Password)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, out)
	case http.MethodDelete:
		userID := r.URL.Query().Get("id")
		if userID == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		if err := s.app.DeleteUser(token, userID); err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, map[string]bool{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	token := tokenFrom(r)
	switch r.Method {
	case http.MethodGet:
		list, err := s.app.ListOrders(token)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, list)
	case http.MethodPost:
		var req models.CreateOrderRequest
		if !readJSON(r, &req) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		out, err := s.app.CreateOrder(token, req)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, out)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		OrderID string `json:"orderId"`
		Status  string `json:"status"`
	}
	if !readJSON(r, &body) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	out, err := s.app.UpdateOrderStatus(tokenFrom(r), body.OrderID, body.Status)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handlePreviewPrice(w http.ResponseWriter, r *http.Request) {
	var q struct {
		TariffID   string  `json:"tariffId"`
		DistanceKm float64 `json:"distanceKm"`
		IdleHours  float64 `json:"idleHours"`
		Urgent     bool    `json:"urgent"`
	}
	if r.Method == http.MethodPost {
		if !readJSON(r, &q) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	} else {
		q.TariffID = r.URL.Query().Get("tariffId")
	}
	price, err := s.app.PreviewPrice(tokenFrom(r), q.TariffID, q.DistanceKm, q.IdleHours, q.Urgent)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, map[string]float64{"price": price})
}

func (s *Server) handleSchedule(w http.ResponseWriter, r *http.Request) {
	list, err := s.app.VehicleSchedule(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, list)
}

func (s *Server) handleWaybills(w http.ResponseWriter, r *http.Request) {
	list, err := s.app.ListWaybills(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, list)
}

func (s *Server) handleOpenShift(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.OpenShiftRequest
	if !readJSON(r, &req) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	out, err := s.app.OpenShift(tokenFrom(r), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handleCloseShift(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.CloseShiftRequest
	if !readJSON(r, &req) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	out, err := s.app.CloseShift(tokenFrom(r), req)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handleShiftStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	out, err := s.app.ShiftStatus(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handleDriversAvailable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	out, err := s.app.ListDriversAvailable(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	out, err := s.app.Dashboard(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handleDriverRating(w http.ResponseWriter, r *http.Request) {
	out, err := s.app.DriverRating(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	out, err := s.app.ListAudit(tokenFrom(r))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, out)
}

type authedHandler func(http.ResponseWriter, *http.Request)

func (s *Server) auth(h authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tokenFrom(r) == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func tokenFrom(r *http.Request) string {
	if t := bearerToken(r); t != "" {
		return t
	}
	return r.URL.Query().Get("token")
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return r.Header.Get("X-Session-Token")
}

func readJSON(r *http.Request, v any) bool {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v) == nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	if ae, ok := err.(*app.AppError); ok {
		w.WriteHeader(http.StatusForbidden)
		writeJSON(w, ae)
		return
	}
	if err.Error() == "unauthorized" || strings.Contains(err.Error(), "invalid credentials") {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	writeJSON(w, map[string]string{"message": err.Error()})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Session-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func spaHandler(staticDir string, fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		path := filepath.Join(staticDir, filepath.Clean("/"+r.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		index := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(index); err == nil {
			r.URL.Path = "/"
			http.ServeFile(w, r, index)
			return
		}
		io.WriteString(w, "frontend not built")
	})
}
