BINARY := tutamen-ask-password

$(BINARY): $(wildcard *.go)
	go build

devel: $(BINARY)
	ln -sf $$PWD/tutamen-ask-password          /usr/lib/systemd/
	ln -sf $$PWD/tutamen-ask-password.service  /usr/lib/systemd/system/
	ln -sf $$PWD/tutamen-ask-password.path     /usr/lib/systemd/system/
	ln -sf $$PWD/sd-tutamen                    /etc/initcpio/install/
	ln -sf $$PWD/sd-networkd                   /etc/initcpio/install/

clean:
	rm -f $(BINARY)

.PHONY: clean devel install
