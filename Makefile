repl-docker:
	docker run --init -it --rm --name chrome --shm-size=1024m --cap-add=SYS_ADMIN --entrypoint=/opt/google/chrome-unstable/chrome yukinying/chrome-headless-browser --headless --disable-gpu --repl https://www.facebook.com
