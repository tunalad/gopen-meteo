PACKAGE_NAME = gopen-meteo
BINARY_NAME = $(PACKAGE_NAME)
INSTALL_DIR = $(HOME)/.local/bin

all: build

build:
	go build

install: build
	install -m 755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	command rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	go clean

lean:
	@echo -e "\e[1;35m💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜\e[0m"
	@echo -e "\e[1;35m💜💜💜💜💜💜💜💜I LOVE LEAN!!!💜💜💜💜💜💜💜💜💜\e[0m"
	@echo -e "\e[1;35m💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜\e[0m"
	@echo -e "\e[1;35m💜I'M ON 'EM BEANS FOR REAL (YEH, YEAH, YEAH)💜\e[0m"
	@echo -e "\e[1;35m💜I'M ON THE LEAN FOR REAL (WHAT? YEAH, YEAH)💜\e[0m"
	@echo -e "\e[1;35m💜I'M ON 'EM BEANS FOR REAL (YEA, YEAH, YEAH)💜\e[0m"
	@echo -e "\e[1;35m💜💜💜I'M ON THE LEAN FOR REAL (YEAH-YEAH)💜💜💜\e[0m"
	@echo -e "\e[1;35m💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜💜\e[0m"
	@echo -e "\e[1;35m💜💜💜💜💜💜💜💜!!LEANEANEAN!!💜💜💜💜💜💜💜💜💜\e[0m"
	@exit 1

.PHONY: all build install uninstall clean lean
