all: 
	rm -f rtty
	go fmt .
	GOOS=linux GOARCH=arm GOARM=6 go build
	@echo
	@echo Copy rtty to target system and then run these commands:
	@echo sudo chown root rtty
	@echo sudo chmod 4755 rtty
clean:
	rm -f rtty
