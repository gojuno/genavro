package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/gojuno/genavro/avro"
	"github.com/mkorolyov/astparser"
)

var (
	inputDir         = flag.String("in", "", "directory with go files to be parsed")
	excludeRegexpStr = flag.String("e", "", "exclude regexp to skip files")
	includeRegexpStr = flag.String("i", "", "include regexp to limit input files")
	outputDir        = flag.String("o", "", "directory for generated avro schemas")
	namespace        = flag.String("n", "", "namespace for generated avro schemas")
)

func main() {
	flag.Parse()

	// load golang sources
	cfg := astparser.Config{InputDir: *inputDir}
	if *excludeRegexpStr != "" {
		cfg.ExcludeRegexp = *excludeRegexpStr
	}
	if *includeRegexpStr != "" {
		cfg.IncludeRegexp = *includeRegexpStr
	}
	sources, err := astparser.Load(cfg)
	if err != nil {
		log.Fatalf("failed to load sources from %s excluding %s: %v", *inputDir, *excludeRegexpStr, err)
	}

	// generate avro protocols
	avroProtocols := avro.Generate(sources, *namespace)

	// save
	for f, r := range avroProtocols {
		filePath := *outputDir + "/" + f + ".avpr"
		bytes, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			log.Fatalf("failed to marshall to file %s generated protocol %+v: %v", f, r, err)
		}
		ioutil.WriteFile(filePath, bytes, 0666)
	}
}
