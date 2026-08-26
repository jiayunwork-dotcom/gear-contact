package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gear-contact/internal/gears"
	"gear-contact/internal/hertz"
	"gear-contact/internal/lewis"
	"gear-contact/internal/mesh"
	"gear-contact/internal/share"
	"gear-contact/internal/slide"
)

const maxBodyBytes = 1 << 20

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
}

type MeshResponse struct {
	CenterDistance float64 `json:"center_distance"`
	BasePitch      float64 `json:"base_pitch"`
	LineLength     float64 `json:"line_length"`
	ContactRatio   float64 `json:"contact_ratio"`
	Clearance      float64 `json:"clearance"`
	Standard       bool    `json:"standard"`
}

type HertzRequest struct {
	Spec json.RawMessage `json:"spec"`
	Ft   float64         `json:"ft"`
	Face float64         `json:"face"`
}

type HertzResponse struct {
	MaxPressure float64 `json:"max_pressure"`
	HalfWidth   float64 `json:"half_width"`
	RhoEq       float64 `json:"rho_eq"`
	PeakP0      float64 `json:"peak_p0"`
	PeakS       float64 `json:"peak_s"`
	PinionBend  float64 `json:"pinion_bend"`
}

func New(staticDir, examplePath string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/example", makeExampleHandler(examplePath))
	mux.HandleFunc("/api/mesh", handleMesh)
	mux.HandleFunc("/api/hertz", handleHertz)
	mux.Handle("/", fileServer(staticDir))
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	writeJSON(w, http.StatusOK, HealthResponse{OK: true, Service: "gear-contact"})
}

func handleMesh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err := gears.FromJSON(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	m, err := mesh.Compute(p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, MeshResponse{
		CenterDistance: m.OperatingCenterDistance,
		BasePitch:      m.BasePitch,
		LineLength:     m.LineOfAction.Length,
		ContactRatio:   m.ContactRatio,
		Clearance:      m.Clearance,
		Standard:       m.Standard,
	})
}

func handleHertz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req HertzRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Ft <= 0 || req.Face <= 0 {
		writeError(w, http.StatusBadRequest, "ft and face must be positive")
		return
	}
	p, err := gears.FromJSON(req.Spec)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	mres, err := mesh.Compute(p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	fn := hertz.NormalFromTangential(req.Ft, mres.WorkingPressureAngle)
	c, err := hertz.AtPitch(p, fn, req.Face)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	samples, err := share.Walk(mres, req.Ft, req.Face, 21)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	peak, _ := share.PeakPressure(samples)
	s1, _, err := lewis.PairBending(p, req.Ft, req.Face)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = slide.Sliding(p, peak.S, 1)
	writeJSON(w, http.StatusOK, HertzResponse{
		MaxPressure: c.MaxPressure,
		HalfWidth:   c.HalfWidth,
		RhoEq:       c.RhoEq,
		PeakP0:      peak.P0,
		PeakS:       peak.S,
		PinionBend:  s1,
	})
}

func makeExampleHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "GET required")
			return
		}
		f, err := os.Open(path)
		if err != nil {
			writeError(w, http.StatusNotFound, "example unavailable: "+err.Error())
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("request body is not valid JSON: %v", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func fileServer(dir string) http.Handler {
	inner := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		inner.ServeHTTP(w, r)
	})
}
