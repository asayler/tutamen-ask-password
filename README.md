# tutament-ask-password

This utility implements the [systemd Password Agent
Specification](http://www.freedesktop.org/wiki/Software/systemd/PasswordAgents).
It can be used in an initrd to unlock encrypted root filesystems using Tutamen.
A password prompt is still available on the console for manual unlocking.

`sd-networkd` and `sd-tutamen` are Archlinux mkinitcpio install hooks.
`sd-networkd` is in this repo simply because it doesn't exist anywhere else.
`sd-tutamen` makes sure that the tutamen specific files are added to the
initrd.

To configure a system that already is set up for LUKS hard disk encryption,
make sure `sd-networkd sd-tutamen` are in the `HOOKS` setting in
`/etc/mkinitcpio.conf` and then

	make
	sudo make devel
	sudo mkinitcpio -P

There are currently some hardcoded values but the system works overall.
