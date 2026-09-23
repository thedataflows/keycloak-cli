package cmd

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// The MCP stdio smoke test drives the real binary over stdin/stdout, so
// TestMain builds it once.
var mcpBinary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "kcc-mcp-smoke-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	mcpBinary = filepath.Join(dir, "keycloak-cli")
	build := exec.Command("go", "build", "-o", mcpBinary, ".")
	build.Dir = ".." // repo root
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("build keycloak-cli binary: " + err.Error())
	}
	os.Exit(m.Run())
}

type rpcResponse struct {
	Result struct {
		ServerInfo struct {
			Name string `json:"name"`
		} `json:"serverInfo"`
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	} `json:"result"`
}

// mcpSpecFlag is the explicit spec flag: SpecPath no longer has a default, so
// the stdio smoke test names the vendored spec.
func mcpSpecFlag(t *testing.T) string {
	t.Helper()
	spec, err := filepath.Abs(filepath.Join("..", "keycloak-oapi", "26.7.4.spec.json"))
	if err != nil {
		t.Fatal(err)
	}
	return spec
}

// Scenario 11: the built binary completes an MCP initialize handshake over
// stdio and lists the five tools.
func TestMcpStdioHandshake(t *testing.T) {
	cmd := exec.Command(mcpBinary, "mcp", "--spec-path", mcpSpecFlag(t))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	reader := bufio.NewReader(stdout)
	readResponse := func() rpcResponse {
		t.Helper()
		type res struct {
			line string
			err  error
		}
		got := make(chan res, 1)
		go func() {
			line, err := reader.ReadString('\n')
			got <- res{line, err}
		}()
		select {
		case r := <-got:
			if r.err != nil {
				t.Fatalf("read rpc response: %v", r.err)
			}
			var out rpcResponse
			if err := json.Unmarshal([]byte(r.line), &out); err != nil {
				t.Fatalf("decode rpc response %q: %v", r.line, err)
			}
			return out
		case <-time.After(30 * time.Second):
			t.Fatal("timed out waiting for rpc response")
			return rpcResponse{}
		}
	}
	send := func(v map[string]any) {
		t.Helper()
		line, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := stdin.Write(append(line, '\n')); err != nil {
			t.Fatalf("write request: %v", err)
		}
	}

	send(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": map[string]any{
		"protocolVersion": "2025-06-18",
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]any{"name": "smoke", "version": "0"},
	}})
	init := readResponse()
	if init.Result.ServerInfo.Name == "" {
		t.Fatalf("initialize returned no serverInfo: %+v", init)
	}

	send(map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"})
	send(map[string]any{"jsonrpc": "2.0", "id": 2, "method": "tools/list"})
	list := readResponse()

	names := make([]string, 0, len(list.Result.Tools))
	for _, tool := range list.Result.Tools {
		names = append(names, tool.Name)
	}
	if len(names) != 5 {
		t.Fatalf("tools/list returned %d tools (%v), want 5", len(names), names)
	}
	_ = stdin.Close()
	_ = cmd.Wait()
}
