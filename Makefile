MAKEFLAGS += --no-print-directory

.SILENT:
.PHONY: src, run, run-devel, target, build, rc, rs, r, b, bc, watch, w, clean, c, install, i, watchexec, we, help, lo 

# don't chang this shit

LO_SRC_DIR = ./src
# LO_TARGET_DIR = $(LO_SRC_DIR)/target
# LO_BIN_DIR = $(LO_TARGET_DIR)/release
# LO_BIN_NAME = "LO"
ORIGIN_DIR = $(shell pwd)

h: help
r: run
b: build
bc: buildcopy
w: watch
we: watchexec
c: clean
i: install
build: b


help: 
	@echo ""
	@echo "| --------- LO Help list --------- |"
	@echo ""
	@echo "Commands:"
	@printf "  run        [r] \tRuns code\n"
	@printf "  build      [b] \tBuilds code\n"
	@printf "  watchexec  [we]\tWatchexec's the code for hot reloads\n"
	@printf "  install    [i] \tInstalls LO via the install script\n"
	@echo ""
	@printf "  help       [h] \tShows this thing\n"


run:
	@echo "| running.......... |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| running code...   |"; \
		go run main.go;


watchexec:
	@echo "| watchexec code    |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| watching code...  |"; \
		watchexec -r "go run main.go";

build:
	@echo "| building code     |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| building code...  |"; \
		go build -ldflags "-s -w" -trimpath -o lo main.go;

buildcopy:
	@echo "| building code     |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| building code...  |"; \
		go build -ldflags "-s -w" -trimpath -o lo main.go; \
		cp lo $(ORIGIN_DIR)/lo;

install:
	@echo "| installing bin... |"
	bash ./install.sh
	@echo "| installed bin 👍  |"


bci: buildcopy install 

