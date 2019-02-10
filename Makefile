
PROG = timeshift

$(PROG) : $(PROG).go
	go build -ldflags "-s -w" $(PROG).go

clean:
	rm -f $(PROG) *~ .??*~

distclean: clean
	rm -rf dist/

