gopher
======

A simple, low maintenance Gopher server. Serves a directory tree with
associated metadata files as a gopher hole.

Usage
-----

    -d string
        The gopher root dir (default ".")
    -n string
        The hostname (default your hostname)
    -p int
        The port to listen to (default 70)
    -w string
        HTTP server address

Metadata
--------

`gopher` reads `.head` files in its directories to produce page headers.
These contain gopher selectors. The selectors only have to be partially
written, and the server will fill missing columns with dummy data. This
can be used to create a header for a file listing, or to produce an
arbitrary Gopher page.

HTTP
----

`gopher` will optionally serve HTTP, given an address in the form
"host:port". This will mirror the gopher content and optionally use HTML
templates:

-   `.template`: The template used for gopher menus
-   `.mdtemplate`: The template used for Markdown documents

