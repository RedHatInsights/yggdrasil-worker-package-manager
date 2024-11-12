package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/subpop/go-log"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/fftoml"
	"github.com/redhatinsights/yggdrasil/ipc"
	"github.com/redhatinsights/yggdrasil/worker"
	"github.com/sgreben/flagvar"
	"github.com/zcalusic/sysinfo"
)

type Message struct {
	Command string `json:"command"`
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

var (
	Version   string
	ConfigDir string
)

var (
	logLevel      = flagvar.Enum{Choices: []string{"error", "warn", "info", "debug", "trace"}}
	allowPatterns = flagvar.Regexps{}
	version       = false
)

func main() {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	fs.Var(
		&logLevel,
		"log-level",
		"log verbosity level (error (default), warn, info, debug, trace)",
	)
	fs.Var(
		&allowPatterns,
		"allow-pattern",
		"regular expression pattern to allow package operations\n(can be specified multiple times)",
	)
	fs.BoolVar(&version, "version", false, "show version info")
	_ = fs.String(
		"config",
		filepath.Join(ConfigDir, "config.toml"),
		"path to `file` containing configuration values (optional)",
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix(), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(fftoml.Parser), ff.WithAllowMissingConfigFile(true)); err != nil {
		log.Fatal(err)
	}

	if version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if logLevel.Value != "" {
		l, err := log.ParseLevel(logLevel.Value)
		if err != nil {
			log.Fatalf("cannot parse log level: %v", err)
		}
		log.SetLevel(l)
	}

	if log.CurrentLevel() >= log.LevelDebug {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	worker, err := worker.NewWorker(
		"package_manager",
		false,
		map[string]string{"version": Version},
		nil,
		dataRx,
		nil,
	)
	if err != nil {
		log.Fatalf("error: cannot create worker: %v", err)
	}

	// Set up a channel to receive the TERM or INT signal over and clean up
	// before quitting.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	if err := worker.Connect(quit); err != nil {
		log.Fatalf("error: cannot connect: %v", err)
	}
}

func dataRx(
	w *worker.Worker,
	addr string,
	id string,
	responseTo string,
	metadata map[string]string,
	data []byte,
) error {
	log.Debugf("received message: %v", id)
	log.Tracef("%v", data)

	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("cannot unmarshal data: %v", err)
	}
	log.Debugf("received command: %v", m)

	pm, err := detectPackageManager()
	if err != nil {
		log.Fatalf("cannot detect package manager: %v", err)
	}

	for _, ch := range []chan []byte{pm.Stdout(), pm.Stderr()} {
		go func(ch chan []byte) {
			for buf := range ch {
				if err := w.EmitEvent(ipc.WorkerEventNameWorking, id, "", map[string]string{"output": strings.TrimRight(string(buf), "\n\x00")}); err != nil {
					log.Errorf("cannot emit event: %v", err)
				}
			}
		}(ch)
	}

	var outb, errb []byte
	var code int
	switch m.Command {
	case "install":
		if !packageAllowed(m.Name) {
			return fmt.Errorf("cannot install %v: does not match an allow pattern", m.Name)
		}
		outb, errb, code, err = pm.Install(m.Name)
	case "remove":
		outb, errb, code, err = pm.Uninstall(m.Name)
	case "enable-repo":
		outb, errb, code, err = pm.EnableRepo(m.Name)
	case "disable-repo":
		outb, errb, code, err = pm.DisableRepo(m.Name)
	case "add-repo":
		outb, errb, code, err = pm.AddRepo(m.Name, m.Content)
	case "remove-repo":
		outb, errb, code, err = pm.RemoveRepo(m.Name)
	default:
		return fmt.Errorf("unknown command: %v", m.Command)
	}

	if err != nil {
		switch e := err.(type) {
		case ExitError:
			log.Errorf("failed running command: %v", e)
		default:
			return fmt.Errorf("cannot run command: %v", e)
		}
	}

	response := map[string]interface{}{
		"code":   code,
		"stdout": string(outb),
		"stderr": string(errb),
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("cannot marshal json: %v", err)
	}

	_, _, _, err = w.Transmit(addr, uuid.New().String(), id, nil, responseData)
	if err != nil {
		return fmt.Errorf("cannot call Transmit: %v", err)
	}

	return nil
}

func packageAllowed(name string) bool {
	for _, r := range allowPatterns.Values {
		if r.Match([]byte(name)) {
			return true
		}
	}
	return false
}

func detectPackageManager() (PackageManager, error) {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	switch si.OS.Vendor {
	case "fedora":
		return &PackageManagerDnf{}, nil
	case "centos", "rhel":
		ver := strings.Split(si.OS.Version, ".")
		if len(ver) == 0 {
			return nil, fmt.Errorf("cannot split version: %v", si.OS.Version)
		}
		major, err := strconv.ParseInt(ver[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse major version component: %w", err)
		}
		if major >= 8 {
			return &PackageManagerDnf{make(chan []byte), make(chan []byte)}, nil
		}
		return &PackageManagerYum{make(chan []byte), make(chan []byte)}, nil
	case "debian", "ubuntu":
		return &PackageManagerApt{make(chan []byte), make(chan []byte)}, nil
	default:
		return nil, fmt.Errorf("unsupported OS: %v", si.OS.Vendor)
	}
}
