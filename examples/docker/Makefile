build:
	docker build . -t headless-chrome

test: build run

run:
	docker run --shm-size 1G --rm -it --cap-add=SYS_ADMIN headless-chrome

repl:
	docker run --shm-size 1G --rm -it --cap-add=SYS_ADMIN -v --entrypoint="/opt/google/chrome-unstable/chrome" headless-chrome --headless --disable-gpu --repl --no-sandbox 'https://ebay.com'
