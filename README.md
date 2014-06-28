# vuvuzela proxy

redirect and inject vuvuzela on pages that are not vuvuzela enabled.
check it at (http://vuvuzela.buildfactory.me)

## Preparing the environment

Prerequisites:

- Git
- rsync
- GNU Make
- [Go](http://golang.org) 1.0.3 or newer

Make sure the Go compiler is installed and `$GOPATH` is set.

Install dependencies, and compile:

	make deps
	make clean
	make all

Generate a self-signed SSL certificate (optional):

	cd ssl
	make

Start Redis if you plan to use it

Edit the config file and run the server. You can either serve it through nginx or straight to the interwebs. 

	vi go-vuvuzela.conf
	./go-vuvuzela

Install, uninstall. Edit Makefile and set PREFIX to the target directory:

	sudo make install
	sudo make uninstall

Allow non-root process to listen on low ports:

	/sbin/setcap 'cap_net_bind_service=+ep' /opt/go-vuvuzela/server

