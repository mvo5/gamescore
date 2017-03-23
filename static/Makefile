
all: deps 

deps:
	./get-jquery.sh

install:
	mkdir -p $(DESTDIR)/static
	install -m644 -t $(DESTDIR)/static *
