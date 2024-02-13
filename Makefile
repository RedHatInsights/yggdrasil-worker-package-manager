VERSION := 0.1.3

.PHONY: dist
dist: yggdrasil-worker-package-manager-$(VERSION).tar.xz

yggdrasil-worker-package-manager-$(VERSION).tar.xz: TMPDIR := $(shell mktemp -d)
yggdrasil-worker-package-manager-$(VERSION).tar.xz:
	go mod vendor
	tar --create \
		--lzma \
		--absolute-names \
		--file $(TMPDIR)/$@ \
		--transform "s+$(PWD)+yggdrasil-worker-package-manager-$(VERSION)+" \
		--exclude $@ \
		--exclude .git \
		--exclude .github \
		--exclude .gitignore \
		--exclude .vscode \
		--exclude builddir \
		--exclude .copr \
		$(PWD)
	mv $(TMPDIR)/$@ .
	rm -rf vendor

yggdrasil-worker-package-manager-$(VERSION).tar.xz.sha256sum: yggdrasil-worker-package-manager-$(VERSION).tar.xz
	shasum -a 256 $^ > $@
