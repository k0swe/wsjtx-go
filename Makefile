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

# TODO: This target can be removed once the package is in Debian stable and Ubuntu stable, 2021-05
leemcloughlin-jdn.deb:
	cd .. && \
	wget http://ftp.debian.org/debian/pool/main/g/golang-github-leemcloughlin-jdn/golang-github-leemcloughlin-jdn-dev_0.0~git20201102.6f88db6-2_all.deb

.PHONY: deb-package
deb-package: deb-orig-tarball leemcloughlin-jdn.deb
	# https://wiki.debian.org/sbuild
	sbuild -d unstable \
      --extra-package=../golang-github-leemcloughlin-jdn-dev_0.0~git20201102.6f88db6-2_all.deb

.PHONY: update-chroots
update-chroots:
	sudo sbuild-update -udcar `ls /srv/chroot/ | grep sbuild`
