init-project:
	go mod init github.com/renegmed/bookstore_oauth-go 

test:
	go test -race ./... 
