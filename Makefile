STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/local/anthropic/anthropic.plugin -tags "$(BUILD_TAGS)" *.go

# Build the plugin as a SQLite loadable extension (steampipe_sqlite_anthropic.so)
# using turbot/steampipe-sqlite with a replace directive pointing at this repo.
sqlite:
	rm -rf build/steampipe-sqlite && mkdir -p build
	git clone --quiet --depth 1 https://github.com/turbot/steampipe-sqlite.git build/steampipe-sqlite
	cd build/steampipe-sqlite && rm -rf work && mkdir -p work && rsync -a --exclude='.git' . work/
	cd build/steampipe-sqlite/work && \
		go run generate/generator.go templates . anthropic "" github.com/huntbase-io/steampipe-plugin-anthropic && \
		echo 'replace github.com/huntbase-io/steampipe-plugin-anthropic => $(CURDIR)' >> go.mod && \
		go mod tidy && \
		$(MAKE) -f out/Makefile build
	cp build/steampipe-sqlite/work/steampipe_sqlite_anthropic.so .
	@echo "Built ./steampipe_sqlite_anthropic.so"

clean:
	rm -rf build steampipe_sqlite_anthropic.so
