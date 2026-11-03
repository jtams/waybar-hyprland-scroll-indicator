package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Client struct {
	Address   string `json:"address"`
	At        [2]int `json:"at"`
	Floating  bool   `json:"floating"`
	Hidden    bool   `json:"hidden"`
	Mapped    bool   `json:"mapped"`
	Workspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workspace"`
}

type ActiveWindow struct {
	Address   string `json:"address"`
	Workspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workspace"`
}

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

func hyprctl(args ...string) ([]byte, error) {
	return exec.Command("hyprctl", args...).Output()
}

func emit() {
	clientsOut, err := hyprctl("clients", "-j")
	if err != nil {
		return
	}
	activeOut, err := hyprctl("activewindow", "-j")
	if err != nil {
		return
	}

	var clients []Client
	var active ActiveWindow

	if err := json.Unmarshal(clientsOut, &clients); err != nil {
		return
	}
	if err := json.Unmarshal(activeOut, &active); err != nil {
		return
	}

	var workspace []Client
	for _, c := range clients {
		if c.Workspace.Name == active.Workspace.Name && !c.Floating && c.Mapped && !c.Hidden {
			workspace = append(workspace, c)
		}
	}

	sort.Slice(workspace, func(i, j int) bool {
		return workspace[i].At[0] < workspace[j].At[0]
	})

	activeIdx := -1
	for i, c := range workspace {
		if c.Address == active.Address {
			activeIdx = i
			break
		}
	}

	total := len(workspace)
	var dots []string
	for i := range workspace {
		if i == activeIdx {
			dots = append(dots, "●")
		} else {
			dots = append(dots, "○")
		}
	}

	text := strings.Join(dots, " ")
	tooltip := fmt.Sprintf("%d / %d", activeIdx+1, total)
	if activeIdx == -1 {
		tooltip = fmt.Sprintf("? / %d", total)
	}

	class := "scrolling"
	if total <= 1 {
		class = "scrolling single"
	}

	out := WaybarOutput{Text: text, Tooltip: tooltip, Class: class}
	b, _ := json.Marshal(out)
	fmt.Println(string(b))
}

func socketPath() string {
	sig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if runtime == "" {
		runtime = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return fmt.Sprintf("%s/hypr/%s/.socket2.sock", runtime, sig)
}

func main() {
	emit()

	conn, err := net.Dial("unix", socketPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to hyprland socket: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	var carry string

	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		data := carry + string(buf[:n])
		carry = ""

		lines := strings.Split(data, "\n")
		for i, line := range lines {
			if i == len(lines)-1 {
				carry = line
				break
			}
			if shouldUpdate(line) {
				emit()
			}
		}
	}
}

func shouldUpdate(event string) bool {
	triggers := []string{
		"activewindow>>",
		"openwindow>>",
		"closewindow>>",
		"movewindow>>",
		"workspace>>",
		"focusedmon>>",
	}
	for _, t := range triggers {
		if strings.HasPrefix(event, t) {
			return true
		}
	}
	return false
}
