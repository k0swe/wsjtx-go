VERSION = $(shell < debian/changelog head -1 | egrep -o "[0-9]+\.[0-9]+\.[0-9]+")

.PHONY: all
all: test

.PHONY: test
test:
	go test ./...
	go vet ./...

.PHONY: deb-orig-tarball
deb-orig-tarball:
	cd .. && tar -cvJf golang-github-k0swe-wsjtx-go_$(VERSION).orig.tar.xz --exclude-vcs --exclude=debian --exclude=.github --exclude=.idea wsjtx-go

.PHONY: deb-tarball
deb-tarball:
	cd .. && tar -cvJf golang-github-k0swe-wsjtx-go_$(VERSION).orig.tar.xz --exclude-vcs wsjtx-go

.PHONY: deb-package
deb-package: deb-tarball
	# https://wiki.debian.org/sbuild
	sbuild -d unstable

.PHONY: update-chroots
update-chroots:
	sudo sbuild-update -udcar `ls /srv/chroot/`
