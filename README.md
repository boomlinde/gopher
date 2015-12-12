gopher
======

A simple, low maintenance Gopher server. Serves a directory tree with
associated metadata files as a gopher hole.

Usage
-----

    -d string
       The gopher root dir (default ".")
    -e int
       The externally visible port (default 70)
    -i int
       The port to listen to (default 7000)
    -n string
       The hostname (defaults to your hostname)

Gopher servers normally listens to connections on a privileged port. The
internal and external flags are used to let the server run unprivileged while
still naming a privileged port for its selectors, so that a privileged program,
e.g. `rinetd`, can forward connections to the unprivileged port.

Metadata
--------

`gopher` reads `.head` files in its directories to produce page headers. These
contain gopher selectors. The selectors only have to be partially written, and
the server will fill missing columns with dummy data. This can be used to
create a header for a file listing, or to produce an arbitrary Gopher page.
