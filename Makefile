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
	@printf "  run        [r] \tRuns the rust code\n"
	@printf "  build      [b] \tBuilds rust code\n"
	@printf "  buildcopy  [bc]\tBuilds & copies binary to root dir\n"
	@printf "  watch      [w] \tCargo watch-es the code for hot reloads\n"
	@printf "  watchexec  [we]\tWatchexec's the code for hot reloads\n"
	@printf "  clean      [c] \tCleans all leftover build targets & others\n"
	@printf "  install    [i] \tInstalls LO via the install script\n"
	@printf "  bci        []  \tBulid copy and install\n"
	@echo ""
	@printf "  help       [h] \tShows this thing\n"




run:
	@echo "| running rust code |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| running code...   |"; \
		cargo run;


watch: 
	@echo "| watching rust code |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| watching code...  |"; \
		cargo watch -x run;


watchexec:
	@echo "| watchexec code    |"
	cd $(LO_SRC_DIR); \
		echo "| entered directory |>" $$PWD; \
		echo "| watching code...  |"; \
		watchexec -r "cargo run";



