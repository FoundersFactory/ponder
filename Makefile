BINARY = ponder
SOURCES = ponder.go

GPG_LIBS = libgpg-error-1.21 libassuan-2.4.2 gpgme-1.6.0

LIB_ROOT = $(CURDIR)/lib
LIB_LOCAL = $(LIB_ROOT)/usr/local

LIB_NAMES = $(foreach lib,$(GPG_LIBS),$(shell echo $(lib) | sed 's/-[0-9.]*$$//'))
LIB_DIRS = $(addprefix $(LIB_ROOT)/,$(LIB_NAMES))
LIB_ARCHIVES = $(addsuffix .a,$(addprefix $(LIB_LOCAL)/lib/,$(LIB_NAMES)))

SRC_DIR = $(CURDIR)/src

export CFLAGS = -I$(LIB_LOCAL)/include
export LDFLAGS = -L$(LIB_LOCAL)/lib

GO_LDFLAGS = --ldflags '-extldflags "-static"'

export CGO_CFLAGS = $(CFLAGS)
export CGO_LDFLAGS = $(LDFLAGS)

export GOPATH = $(CURDIR)
export GOBIN = $(CURDIR)/bin

.PHONY: build clean libs

build: $(GOBIN)/$(BINARY)

libs: $(LIB_ARCHIVES)

src: $(SRC_DIR)

clean:
	rm -fr $(BINARY) $(SRC_DIR) $(LIB_ROOT)

$(GOBIN)/$(BINARY): $(SOURCES) $(SRC_DIR)
	go install $(GO_LDFLAGS) $<

$(SRC_DIR): $(SOURCES)
	go get -d .
	touch $@

$(LIB_ROOT):
	mkdir -p $@

$(LIB_DIRS): $(LIB_ROOT)/%: $(LIB_ROOT)
	mkdir -p $@
	curl -s https://www.gnupg.org/ftp/gcrypt/$*/$(filter $*%,$(GPG_LIBS)).tar.bz2 |\
		tar --extract --bzip2 --directory=$@ --strip-components=1

$(LIB_ARCHIVES): $(LIB_LOCAL)/lib/%.a: $(LIB_ROOT)/%
	cd $< && ./configure \
		--with-libgpg-error-prefix=$(LIB_LOCAL)\
		--with-libassuan-prefix=$(LIB_LOCAL)\
		--enable-static
	make -C $< install DESTDIR=$(LIB_ROOT)
