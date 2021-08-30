pb:
	protoc -I=std/serialization --go_out=std book.proto

cert:
	openssl req -x509 -newkey rsa:4096 -keyout std/http_server/server.key -out std/http_server/server.crt -days 365 -nodes -subj '/CN=localhost'