package version

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	rawVersions string

	dependencies string

	BuildedAt    string
	Builder      string
	Tag          string
	TagShort     string
	TagCreatedAt string
	CommitHash   string
	Dependencies = make(map[string]string)

	// this variable set by ms_core to avoid import cycle
	Service string
)

var versionFlag = pflag.BoolP("version", "v", false, "Shows versions")

//Init parses  versions of dependencies provided by -ldflags
func Init() {
	if len(rawVersions) == 0 && len(BuildedAt) == 0 {
		fmt.Println("No VERSION info provided.")
		if *versionFlag {
			os.Exit(0)
		}
		return
	}
	oldInit()

	BuildedAt = unifyTime(BuildedAt)
	TagShort, TagCreatedAt = parseTag(Tag)

	parsed := strings.Split(dependencies, ";")
	for _, dep := range parsed {
		res := strings.Split(dep, "=")
		if len(res) == 2 {
			Dependencies[res[0]] = res[1]
		}
	}
	if *versionFlag {
		fmt.Printf("%s\n", versionInfo())
		os.Exit(0)
	}
}

func versionInfo() string {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "Service: %s\n", Service)
	fmt.Fprintf(buffer, "Builder: %s, BuildedAt: %s\n", Builder, BuildedAt)
	fmt.Fprintf(buffer, "CommitHash: %s, Tag: %s, CreatedAt: %s (%s)\n", CommitHash, TagShort, TagCreatedAt, Tag)
	fmt.Fprintf(buffer, "Dependencies:\n")
	for s, v := range Dependencies {
		fmt.Fprintf(buffer, "\t%s: %s\n", s, v)
	}
	return buffer.String()
}

func oldInit() {
	parsedVersions := strings.Split(rawVersions, ";")
	for _, ver := range parsedVersions {
		res := strings.Split(ver, "=")
		if len(res) == 2 {
			key := res[0]
			value := res[1]
			switch key {
			case "date":
				BuildedAt = value
			case "compiler":
				Builder = value
			default:
				Dependencies[key] = value
			}
		}
	}
}

// Versions returns versions of dependencies
func Versions() map[string]string {
	return Dependencies
}

type ServiceVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServiceMeta struct {
	ServiceVersion
	BuildDate    string           `json:"built"`
	Compiler     string           `json:"compiler"`
	Dependencies []ServiceVersion `json:"dependencies"`
}

type VersionRequest struct {
	Subject string `json:"subject"`
}

type VersionResponse struct {
	ServiceMeta
}

// PrettyVersions returns versions of dependencies in structured form
func PrettyVersions() ServiceMeta {
	serviceMeta := ServiceMeta{}
	serviceMeta.Name = Service
	serviceMeta.Dependencies = []ServiceVersion{}
	for key, val := range Dependencies {
		serviceMeta.Dependencies = append(serviceMeta.Dependencies, ServiceVersion{Name: key, Version: val})
	}
	return serviceMeta
}

func parseTag(raw string) (string, string) {
	tag := "<unknown>"
	ts := "<unknown>"

	res := strings.Split(raw, "+")
	if len(res) == 2 {
		tag = res[0]
		ts = unifyTime(res[1])
	}

	return tag, ts
}

func unifyTime(ts string) string {
	result := ts
	if t, err := time.Parse("20060102_150405Z", ts); err == nil {
		result = t.Format(time.RFC3339)
	}
	return result
}
