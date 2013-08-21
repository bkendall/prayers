all: server javascript

server: main.go
	go get

javascript: js_src/*.coffee
	coffee --compile --output js_lib js_src/*.coffee

clean:
	rm -rf js_lib/*
	rm -rf ../../deploy
	rm -rf ../../pkg
	rm -rf ../../bin

deploy:
	mkdir -p ../../deploy
	cp ../../bin/server ../../deploy
	cp -r templates/ ../../deploy
	cp -r js_lib/ ../../deploy
