package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~spc/go-log"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/fftoml"
	"github.com/redhatinsights/yggdrasil"
	pb "github.com/redhatinsights/yggdrasil/protocol"
	"github.com/sgreben/flagvar"
	"github.com/zcalusic/sysinfo"
	"google.golang.org/grpc"
)

func main() {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	var (
		socketAddr    = ""
		logLevel      = flagvar.Enum{Choices: []string{"error", "warn", "info", "debug", "trace"}}
		allowPatterns = flagvar.Regexps{}
	)

	fs.StringVar(&socketAddr, "socket-addr", "", "dispatcher socket address")
	fs.Var(&logLevel, "log-level", "log verbosity level (error (default), warn, info, debug, trace)")
	fs.Var(&allowPatterns, "allow-pattern", "regular expression pattern to allow package operations\n(can be specified multiple times)")
	_ = fs.String("config", filepath.Join(yggdrasil.SysconfDir, yggdrasil.LongName, "workers", fs.Name()+".toml"), "path to `file` containing configuration values (optional)")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("YGG"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(fftoml.Parser)); err != nil {
		log.Fatal(err)
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

	// Dial the dispatcher on its well-known address.
	conn, err := grpc.Dial(socketAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a dispatcher client
	c := pb.NewDispatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Register as a handler of the "package-manager" type.
	r, err := c.Register(ctx, &pb.RegistrationRequest{Handler: "package-manager", Pid: int64(os.Getpid())})
	if err != nil {
		log.Fatal(err)
	}
	if !r.GetRegistered() {
		log.Fatalf("handler registration failed: %v", err)
	}

	// Listen on the provided socket address.
	l, err := net.Listen("unix", r.GetAddress())
	if err != nil {
		log.Fatal(err)
	}

	pm, err := detectPackageManager()
	if err != nil {
		log.Fatalf("cannot detect package manager: %v", err)
	}

	// Register as a Worker service with gRPC and start accepting connections.
	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, &Server{dispatchAddr: socketAddr, allowPatterns: allowPatterns.Values, pm: pm})
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
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
			return &PackageManagerDnf{}, nil
		}
		return &PackageManagerYum{}, nil
	case "debian", "ubuntu":
		return &PackageManagerApt{}, nil
	default:
		return nil, fmt.Errorf("unsupported OS: %v", si.OS.Vendor)
	}
}
