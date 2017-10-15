BUILD_DIR := build

.PHONY: client clean dev

client:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/chat

dev: client
	./build/chat -addr=:8888

clean:
	rm -rf $(BUILD_DIR)
