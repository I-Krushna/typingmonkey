package driver

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/tebeka/selenium"
)

type GeckoDriver struct {
	cmd     *exec.Cmd
	port    int
	wd      selenium.WebDriver
	firefox string
}

func New(port int) (*GeckoDriver, error) {
	firefox, err := findFirefox()
	if err != nil {
		return nil, err
	}

	return &GeckoDriver{
		port:    port,
		firefox: firefox,
	}, nil
}

func (g *GeckoDriver) Start() error {
	g.cmd = exec.Command(".\\geckodriver.exe", "--port", fmt.Sprintf("%d", g.port))
	if err := g.cmd.Start(); err != nil {
		return fmt.Errorf("oopsies!: failed to start geckodriver: %w", err)
	}

	time.Sleep(2 * time.Second)
	return nil
}

func (g *GeckoDriver) Connect() error {
	caps := selenium.Capabilities{
		"browserName": "firefox",
		"moz:firefoxOptions": map[string]interface{}{
			"binary": g.firefox,
		},
	}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", g.port))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	g.wd = wd
	return nil
}

func (g *GeckoDriver) WebDriver() selenium.WebDriver {
	return g.wd
}

func (g *GeckoDriver) Stop() {
	if g.wd != nil {
		g.wd.Quit()
	}
	if g.cmd != nil && g.cmd.Process != nil {
		g.cmd.Process.Kill()
	}
}

func findFirefox() (string, error) {
	paths := []string{
		"C:\\Program Files\\Mozilla Firefox\\firefox.exe",
		"C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe",
		os.Getenv("PROGRAMFILES") + "\\Mozilla Firefox\\firefox.exe",
		os.Getenv("PROGRAMFILES(X86)") + "\\Mozilla Firefox\\firefox.exe",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("oopsies!: firefox not found")
}
