PROJECT_NAME = wd
LOGS_DIR = ./logs

.PHONY: clean

build : logs/
	go build -o $(PROJECT_NAME) main.go

logs/ :
	mkdir $(LOGS_DIR)

clean :
	rm -rf $(PROJECT_NAME)
	rm -rf $(LOGS_DIR)
	rm -rf *.log

