build:
	go build -i -o bin/genavro cmd/genavro/main.go

clean:
	rm -f bin/genavro
	find ./api -maxdepth 1 -name "*.go" | xargs rm
	rm -f api/avro/*


gen: clean build
	cp ../../../junolab.net/ms_data_bridge/ms_data_bridge/* api/
	bin/genavro -in api -o api/avro -n net.junolab.dwh -e "(easyjson.go|client|data_bridge|event_types|event_validator|message_types|metrics_counter|raw_message|requests|test_beqa|transport|validators)"