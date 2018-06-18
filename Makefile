deps:
	dep ensure
install:
	gx install
	dep ensure
	(cd vendor/github.com/Bit-Nation/panthalassa/ && gx install)
	(cd vendor/github.com/Bit-Nation/panthalassa/ && gx-go rw)
	gx-go rw
deps_hack:
	gx-go rw
deps_hack_revert:
	gx-go uw