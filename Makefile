export GOROOT:=$(HOME)/go
export GOARCH:=amd64
export GOOS:=linux
export PATH:=$(PATH):$(HOME)/bin

BUILDDIR=build
SRCDIR=src
OBJ = gogo.6

all: help

gogo: $(OBJ)
	6l -o gogo -L$(BUILDDIR)/ $(BUILDDIR)/gogo.6

%.6: $(SRCDIR)/%.go
	mkdir build
	6g -o $(BUILDDIR)/$@ $<

help:
	@echo "Welcome to the GoGo builder"
	@echo "Configured to use the following go environment:"
	@echo "  GOROOT... $(GOROOT)"
	@echo "  GOARCH... $(GOARCH)"
	@echo "  GOOS... $(GOOS)"
	@echo "Use \`make gogo\` to build the compiler using 6g and 6l"

clean:
	rm -rf build
	rm -f gogo

