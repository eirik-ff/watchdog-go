PROJECT_NAME = wd
LOGS_DIR = ./logs

.PHONY: clean

build : logs/
	go build -o $(PROJECT_NAME) main.go

logs/ :
	mkdir $(LOGS_DIR)

clean :
	rm -rf $(PROJECT_NAME) 2> /dev/null
	rm -rf $(LOGS_DIR) 2> /dev/null
	rm -rf *.log 2> /dev/null

