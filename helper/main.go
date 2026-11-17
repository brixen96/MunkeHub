// MunkeHub helper — watches the WoW SavedVariables file written by the
// MunkeHub addon and POSTs the character stats to the MunkeHub API.
//
// Two payloads:
//
//	syncJSON  — own character stats        -> POST /api/v1/characters
//	relayJSON — guild members received in-game -> POST /api/v1/relay
//
// NOTE: WoW names the SavedVariables file after the ADDON FOLDER
// (MunkeHub.lua), not the ## SavedVariables: variable (MunkeHubDB).
//
// WoW addons cannot do HTTP, so the addon stores stats in SavedVariables
// (written by the client on /reload or logout) and this helper bridges
// that file to the API.
//
// Usage:
//
//	munkehub-helper.exe -token <your-personal-token> [-api http://server:4711] [-file <path>] [-once]
//	munkehub-helper.exe -install     start silently with Windows
//	munkehub-helper.exe -uninstall   remove auto-start
//
// Self-update: at startup the helper asks GET /api/v1/version; if the
// server has a newer exe it downloads it, verifies the SHA256, swaps the
// file (Windows allows renaming a running exe) and restarts itself.
package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// version MUST match helperVersion in api/main.go.
const version = "0.6.0"

// Matches ["syncJSON"]="..." or syncJSON="..." in the SavedVariables file.
var syncRe = regexp.MustCompile(`\[?"syncJSON"?\]\s*=\s*"((?:\\.|[^"\\])*)"`)

// Matches ["relayJSON"]="..." (array of relayed guild members).
var relayRe = regexp.MustCompile(`\[?"relayJSON"?\]\s*=\s*"((?:\\.|[^"\\])*)"`)

// Matches ["detailJSON"]="..." (spec/zone/gear detail for expandable rows).
var detailRe = regexp.MustCompile(`\[?"detailJSON"?\]\s*=\s*"((?:\\.|[^"\\])*)"`)

const defaultGlob = `C:\Program Files (x86)\World of Warcraft\_forever_\WTF\Account\*\SavedVariables\MunkeHub.lua`

func unescapeLua(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case '\\':
			b.WriteByte('\\')
		case '"':
			b.WriteByte('"')
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case 'a':
			b.WriteByte(7)
		case 'b':
			b.WriteByte(8)
		case 'f':
			b.WriteByte(12)
		case 'v':
			b.WriteByte(11)
		default:
			if s[i] >= '0' && s[i] <= '9' {
				j := i
				for j < len(s) && j < i+3 && s[j] >= '0' && s[j] <= '9' {
					j++
				}
				n, _ := strconv.Atoi(s[i:j])
				b.WriteByte(byte(n))
				i = j - 1
			} else {
				b.WriteByte(s[i])
			}
		}
	}
	return b.String()
}

func extractField(content []byte, re *regexp.Regexp, field string) (string, error) {
	m := re.FindSubmatch(content)
	if m == nil {
		return "", fmt.Errorf("%s not found in file (has the addon synced yet? /reload in game first)", field)
	}
	return unescapeLua(string(m[1])), nil
}

// --- self-update -----------------------------------------------------------

type versionInfo struct {
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
	Size    int    `json:"size"`
}

// newerVersion reports whether remote is a higher semver-ish version.
func newerVersion(remote, local string) bool {
	rp := strings.Split(remote, ".")
	lp := strings.Split(local, ".")
	for i := 0; i < len(rp) && i < len(lp); i++ {
		r, _ := strconv.Atoi(rp[i])
		l, _ := strconv.Atoi(lp[i])
		if r != l {
			return r > l
		}
	}
	return len(rp) > len(lp)
}

// checkUpdate downloads a newer helper if available and restarts the
// process. Returns normally if already up to date (or on any error —
// an update problem must never stop the sync loop).
func checkUpdate(apiURL string) {
	client := &http.Client{Timeout: 15 * time.Second}

	res, err := client.Get(strings.TrimRight(apiURL, "/") + "/api/v1/version")
	if err != nil {
		logf("update check skipped: %v", err)
		return
	}
	defer res.Body.Close()
	var vi versionInfo
	if err := json.NewDecoder(io.LimitReader(res.Body, 4096)).Decode(&vi); err != nil {
		logf("update check skipped: %v", err)
		return
	}
	if !newerVersion(vi.Version, version) {
		return // up to date
	}

	logf("update available: %s -> %s, downloading…", version, vi.Version)
	res2, err := client.Get(strings.TrimRight(apiURL, "/") + "/api/v1/helper.exe")
	if err != nil {
		logf("update download failed: %v", err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusOK {
		logf("update download failed: HTTP %d", res2.StatusCode)
		return
	}
	data, err := io.ReadAll(io.LimitReader(res2.Body, 32<<20))
	if err != nil {
		logf("update download failed: %v", err)
		return
	}

	sum := sha256.Sum256(data)
	if hex.EncodeToString(sum[:]) != strings.ToLower(vi.SHA256) {
		logf("update rejected: checksum mismatch")
		return
	}
	if len(data) != vi.Size {
		logf("update rejected: size mismatch")
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		logf("update failed: %v", err)
		return
	}
	newPath := exePath + ".new"
	if err := os.WriteFile(newPath, data, 0755); err != nil {
		logf("update failed: %v", err)
		return
	}
	// Windows allows renaming a running exe; old file becomes .old and is
	// cleaned up on next update.
	oldPath := exePath + ".old"
	_ = os.Remove(oldPath)
	if err := os.Rename(exePath, oldPath); err != nil {
		logf("update failed: %v", err)
		_ = os.Remove(newPath)
		return
	}
	if err := os.Rename(newPath, exePath); err != nil {
		logf("update failed: %v", err)
		_ = os.Rename(oldPath, exePath) // roll back
		return
	}

	logf("updated to %s — restarting", vi.Version)
	cmd := exec.Command(exePath)
	cmd.Args = os.Args
	if silent {
		cmd.Args = append(cmd.Args, "-silent")
	}
	cmd.Dir = filepath.Dir(exePath)
	_ = cmd.Start()
	os.Exit(0)
}

// --- silent mode: log to file instead of console ---------------------------
var (
	silent  bool
	logFile *os.File
)

func logf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if silent {
		if logFile != nil {
			fmt.Fprintf(logFile, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
		}
		return
	}
	fmt.Println(msg)
}

// --- auto-start with Windows (HKCU Run key, no admin needed) --------------
const runKey = `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
const runName = "MunkeHubHelper"

func installSelf(exePath string) error {
	dir := filepath.Dir(exePath)
	vbsPath := filepath.Join(dir, "munkehub-helper-hidden.vbs")
	// VBS launches the exe with NO console window
	vbs := fmt.Sprintf("CreateObject(\"WScript.Shell\").Run \"\"\"%s\"\" -silent\", 0, False", exePath)
	if err := os.WriteFile(vbsPath, []byte(vbs), 0644); err != nil {
		return err
	}
	out, err := exec.Command("reg", "add", runKey, "/v", runName, "/t", "REG_SZ",
		"/d", fmt.Sprintf("wscript.exe \"%s\"", vbsPath), "/f").CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg add: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func uninstallSelf(exePath string) error {
	out, err := exec.Command("reg", "delete", runKey, "/v", runName, "/f").CombinedOutput()
	if err != nil && !strings.Contains(string(out), "unable to find") {
		return fmt.Errorf("reg delete: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	vbsPath := filepath.Join(filepath.Dir(exePath), "munkehub-helper-hidden.vbs")
	_ = os.Remove(vbsPath)
	return nil
}

// fileConfig is written by scripts/mint-token.ps1 next to the helper exe.
type fileConfig struct {
	API   string `json:"api"`
	Token string `json:"token"`
}

func loadConfig() fileConfig {
	var cfg fileConfig
	paths := []string{"config.json"}
	if exe, err := os.Executable(); err == nil {
		paths = append([]string{filepath.Join(filepath.Dir(exe), "config.json")}, paths...)
	}
	for _, p := range paths {
		if b, err := os.ReadFile(p); err == nil {
			if json.Unmarshal(b, &cfg) == nil && (cfg.API != "" || cfg.Token != "") {
				return cfg
			}
		}
	}
	return cfg
}

func findFile(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	matches, err := filepath.Glob(defaultGlob)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("MunkeHub.lua not found yet — log into WoW with the MunkeHub addon enabled, then /reload (or log out) so WoW writes the file, then run this again")
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("multiple WoW accounts found, use -file:\n  %s", strings.Join(matches, "\n  "))
	}
	return matches[0], nil
}

func post(apiURL, token, path, payload string) error {
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(apiURL, "/")+path,
		bytes.NewReader([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
	if res.StatusCode >= 300 {
		return fmt.Errorf("API %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func main() {
	token := flag.String("token", "", "personal token from the MunkeHub admin (required)")
	apiURL := flag.String("api", "http://localhost:4711", "MunkeHub API base URL")
	file := flag.String("file", "", "path to MunkeHub.lua SavedVariables (default: auto-detect)")
	interval := flag.Duration("interval", 10*time.Second, "how often to check the file")
	once := flag.Bool("once", false, "sync once and exit (for testing)")
	install := flag.Bool("install", false, "start automatically with Windows (hidden, no window)")
	uninstall := flag.Bool("uninstall", false, "remove automatic startup")
	silentFlag := flag.Bool("silent", false, "run without console output (log to helper.log)")
	noUpdate := flag.Bool("noupdate", false, "skip the self-update check at startup")
	flag.Parse()

	silent = *silentFlag
	if silent {
		exe, _ := os.Executable()
		logFile, _ = os.OpenFile(filepath.Join(filepath.Dir(exe), "helper.log"),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		defer func() {
			if logFile != nil {
				logFile.Close()
			}
		}()
	}

	exePath, _ := os.Executable()
	if *install {
		if err := installSelf(exePath); err != nil {
			fmt.Println("install failed:", err)
			os.Exit(1)
		}
		fmt.Println("OK - MunkeHub helper now starts automatically with Windows (hidden).")
		fmt.Println("It runs in the background from now on - nothing to remember.")
		return
	}
	if *uninstall {
		if err := uninstallSelf(exePath); err != nil {
			fmt.Println("uninstall failed:", err)
			os.Exit(1)
		}
		fmt.Println("OK - automatic startup removed.")
		return
	}

	cfg := loadConfig()
	if *token == "" {
		*token = cfg.Token
	}
	tokenFromFlag, apiFromFlag := false, false
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "token":
			tokenFromFlag = true
		case "api":
			apiFromFlag = true
		}
	})
	if !apiFromFlag && cfg.API != "" {
		*apiURL = cfg.API
	}

	// Self-update FIRST — before any token prompt or file watching, so an
	// outdated helper always updates even if it can't sync yet. The update
	// endpoints are public (no token needed). Skipped with -once/-noupdate.
	if !*once && !*noUpdate {
		checkUpdate(*apiURL)
	}

	if *token == "" && silent {
		// Never block a background run waiting for input
		logf("no token configured - run munkehub-helper.exe once in a normal window first")
		return
	}
	if *token == "" {
		// First run (generic download): ask once, save for next time
		fmt.Print("No token found. Paste your personal MunkeHub token (ask Emil): ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		*token = strings.TrimSpace(line)
		if *token == "" {
			fmt.Fprintln(os.Stderr, "error: no token provided")
			os.Exit(1)
		}
		cfg.Token = *token
		if cfg.API == "" {
			cfg.API = *apiURL
		}
		cfgPath := "config.json"
		if exe, err := os.Executable(); err == nil {
			cfgPath = filepath.Join(filepath.Dir(exe), "config.json")
		}
		if b, err := json.MarshalIndent(cfg, "", "  "); err == nil {
			_ = os.WriteFile(cfgPath, b, 0600)
		}
		fmt.Println("Token saved — you won't be asked again.")
	}

	path, err := findFile(*file)
	if err != nil {
		if *once {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		// Set-and-forget: wait for WoW to write the file instead of exiting
		logf("waiting for SavedVariables… %v", err)
		for {
			time.Sleep(*interval)
			if p, ferr := findFile(*file); ferr == nil {
				path = p
				break
			}
		}
	}
	logf("watching: %s", path)
	logf("API:      %s", *apiURL)
	logf("version:  %s", version)
	if tokenFromFlag {
		logf("token:    from -token flag")
	} else {
		logf("token:    from config.json")
	}

	var lastMod time.Time
	var lastSize int64

	sync := func() {
		info, err := os.Stat(path)
		if err != nil {
			logf("stat error: %v", err)
			return
		}
		if info.ModTime().Equal(lastMod) && info.Size() == lastSize {
			return // unchanged since last sync
		}
		content, err := os.ReadFile(path)
		if err != nil {
			logf("read error: %v", err)
			return
		}
		// 1. Own stats
		if payload, err := extractField(content, syncRe, "syncJSON"); err == nil && strings.TrimSpace(payload) != "" {
			if err := post(*apiURL, *token, "/api/v1/characters", payload); err != nil {
				logf("post error: %v", err)
				return // don't update lastMod — retry next tick
			}
		} else if err != nil {
			logf("skip: %v", err)
		}

		// 2. Relayed guild members (received in-game via CHAT_MSG_ADDON)
		if relay, err := extractField(content, relayRe, "relayJSON"); err == nil && strings.TrimSpace(relay) != "" {
			if err := post(*apiURL, *token, "/api/v1/relay", relay); err != nil {
				logf("relay error: %v", err)
				return
			}
		}

		// 3. Own character detail (spec/zone/gear) — helper path only
		if detail, err := extractField(content, detailRe, "detailJSON"); err == nil && strings.TrimSpace(detail) != "" {
			if err := post(*apiURL, *token, "/api/v1/characters/detail", detail); err != nil {
				logf("detail error: %v", err)
				return
			}
		}
		lastMod, lastSize = info.ModTime(), info.Size()
		logf("synced at %s", time.Now().Format("15:04:05"))
	}

	sync()
	if *once {
		return
	}
	for range time.Tick(*interval) {
		sync()
	}
}
