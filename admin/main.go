// MunkeHub Admin — a small local web app for Emil to manage guild tokens.
//
// Runs on 127.0.0.1 only (never exposed), opens the browser, and talks to
// the LAN API (:4711) with the admin token.
//
// Usage:
//
//	munkehub-admin.exe [-api http://brixmedia.local:4711] [-port 4713]
//
// First run asks for the admin token once; it is saved to config.json
// (gitignored) next to the exe.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type config struct {
	API        string `json:"api"`
	AdminToken string `json:"adminToken"`
}

var cfg config

func configPath() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "config.json")
	}
	return "config.json"
}

func loadConfig() {
	b, err := os.ReadFile(configPath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &cfg)
}

func saveConfig() error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), b, 0600)
}

// callAPI performs an admin request against the real MunkeHub API.
func callAPI(method, path, body string) (int, string, error) {
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, strings.TrimRight(cfg.API, "/")+path, rdr)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	return res.StatusCode, string(b), nil
}

// --- handlers ---------------------------------------------------------------

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	status, body, err := callAPI("GET", "/api/v1/tokens", "")
	tokens := "[]"
	if err == nil && status == 200 {
		tokens = body
	}
	apiErr := ""
	if err != nil {
		apiErr = err.Error()
	} else if status != 200 {
		apiErr = fmt.Sprintf("API %d: %s", status, body)
	}
	page.Execute(w, map[string]any{
		"API":        cfg.API,
		"HasToken":   cfg.AdminToken != "",
		"Tokens":     template.JS(tokens),
		"APIError":   apiErr,
		"Configured": cfg.API != "" && cfg.AdminToken != "",
	})
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	var in struct {
		API        string `json:"api"`
		AdminToken string `json:"adminToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cfg.API = strings.TrimRight(strings.TrimSpace(in.API), "/")
	cfg.AdminToken = strings.TrimSpace(in.AdminToken)
	if err := saveConfig(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// verify it works
	status, body, err := callAPI("GET", "/api/v1/tokens", "")
	if err != nil || status != 200 {
		msg := fmt.Sprintf("saved, but API check failed: %d %s", status, body)
		if err != nil {
			msg = "saved, but API unreachable: " + err.Error()
		}
		w.WriteHeader(502)
		w.Write([]byte(msg))
		return
	}
	w.Write([]byte("ok"))
}

func handleMint(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Label string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || strings.TrimSpace(in.Label) == "" {
		http.Error(w, "label required", 400)
		return
	}
	status, body, err := callAPI("POST", "/api/v1/tokens",
		fmt.Sprintf(`{"label":%q}`, strings.TrimSpace(in.Label)))
	if err != nil {
		http.Error(w, "API unreachable: "+err.Error(), 502)
		return
	}
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func handleRevoke(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/revoke/")
	if id == "" || id == r.URL.Path {
		http.Error(w, "id required", 400)
		return
	}
	status, body, err := callAPI("DELETE", "/api/v1/tokens/"+id, "")
	if err != nil {
		http.Error(w, "API unreachable: "+err.Error(), 502)
		return
	}
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func main() {
	api := flag.String("api", "http://brixmedia.local:4711", "MunkeHub API (LAN)")
	port := flag.String("port", "4713", "local port for the admin UI")
	noOpen := flag.Bool("no-open", false, "don't open the browser")
	flag.Parse()

	loadConfig()
	if cfg.API == "" {
		cfg.API = strings.TrimRight(*api, "/")
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/mint", handleMint)
	http.HandleFunc("/revoke/", handleRevoke)

	url := "http://127.0.0.1:" + *port + "/"
	fmt.Println("MunkeHub Admin running at", url)
	if !*noOpen {
		_ = exec.Command("cmd", "/c", "start", "", url).Start()
	}
	// 127.0.0.1 only — never exposed beyond this machine
	if err := http.ListenAndServe("127.0.0.1:"+*port, nil); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
