# makefile has a very stupid relation with tabs,
# all actions of every rule are identified by tabs ......
# and No 4 spaces don't make a tab,
# only a tab makes a tab...
#
# to check I use the command 'cat -e -t -v makefile_name'
# It shows the presence of tabs with ^I and line endings with $
# both are vital to ensure that dependencies end properly and tabs mark the action for the rules so that they are easily identifiable to the make utility.....

APP_DIR = app/chat
BUILD_DIR = build/chat

.PHONY: chat

all: chat

chat:
	@echo "Generate directory 'chat' into build..."
	@if [ ! -d $(BUILD_DIR) ] ; then \
	    mkdir $(BUILD_DIR); \
	fi
	@echo "Generate directory 'logs' into $(BUILD_DIR)..."
	@if [ ! -d $(BUILD_DIR)/logs ] ; then \
        mkdir $(BUILD_DIR)/logs; \
    fi

	@echo "Copy config file 'app.toml' into $(BUILD_DIR)..."
	cp $(APP_DIR)/cmd/app.toml $(BUILD_DIR)/app.toml

	@echo "Copy sh file 'go-chat-main.sh' into $(BUILD_DIR)..."
	cp $(APP_DIR)/cmd/go-chat-main.sh $(BUILD_DIR)/go-chat-main.sh
	chmod +x $(BUILD_DIR)/go-chat-main.sh

	echo "Compiling for mac, linux, win platform"
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/main-darwin $(APP_DIR)/cmd/main.go
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/main-linux $(APP_DIR)/cmd/main.go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/main-win $(APP_DIR)/cmd/main.go
	chmod +x $(BUILD_DIR)/main-darwin
	chmod +x $(BUILD_DIR)/main-linux
	chmod +x $(BUILD_DIR)/main-win

	@echo "Make success."
	@echo "You can excute 'cd $(BUILD_DIR)'"
	@echo "Run this service with: './main-darwin' or './go-chat-main.sh start'"
	@echo "Also you can move $(BUILD_DIR) to any place to run this service"
