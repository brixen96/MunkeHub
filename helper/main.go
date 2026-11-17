// MunkeHub helper — watches the WoW SavedVariables file written by the
// MunkeHub addon and POSTs the character stats to the MunkeHub API.
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
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Matches ["syncJSON"]="..." or syncJSON="..." in the SavedVariables file.
var syncRe = regexp.MustCompile(`\[?"syncJSON"?\]\s*=\s*"((?:\\.|[^"\\])*)"`)

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

func extractPayload(content []byte) (string, error) {
	m := syncRe.FindSubmatch(content)
	if m == nil {
		return "", fmt.Errorf("syncJSON not found in file (has the addon synced yet? /reload in game first)")
	}
	return unescapeLua(string(m[1])), nil
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

func post(apiURL, token, payload string) error {
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(apiURL, "/")+"/api/v1/characters",
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
	flag.Parse()

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
	if *token == "" {
		fmt.Fprintln(os.Stderr, "error: no token — run scripts\\mint-token.ps1 or pass -token")
		flag.Usage()
		os.Exit(1)
	}

	path, err := findFile(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("watching:", path)
	fmt.Println("API:     ", *apiURL)
	if tokenFromFlag {
		fmt.Println("token:   from -token flag")
	} else {
		fmt.Println("token:   from config.json")
	}

	var lastMod time.Time
	var lastSize int64

	sync := func() {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("stat error:", err)
			return
		}
		if info.ModTime().Equal(lastMod) && info.Size() == lastSize {
			return // unchanged since last sync
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		payload, err := extractPayload(content)
		if err != nil {
			fmt.Println("skip:", err)
			return
		}
		if strings.TrimSpace(payload) == "" {
			fmt.Println("skip: syncJSON is empty — run /mh sync in game, then /reload")
			return
		}
		if err := post(*apiURL, *token, payload); err != nil {
			fmt.Println("post error:", err)
			return // don't update lastMod — retry next tick
		}
		lastMod, lastSize = info.ModTime(), info.Size()
		fmt.Printf("synced at %s\n", time.Now().Format("15:04:05"))
	}

	sync()
	if *once {
		return
	}
	for range time.Tick(*interval) {
		sync()
	}
}