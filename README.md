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
can be used to create a header for a file listing.

It also reads `.dir` files. These are like `.head` files in that they
contain partial selectors. However, when present, they will replace the
default directory listing altogether. Instead, in a `.dir` file, a
selector with the item type "{" and display string "DIR}" will be
replaced by a listing of the current directory.

File names starting with "." won't be listed in the default listing, but
will still be accessible if selected by a client.

HTTP
----

`gopher` will optionally serve HTTP, given an address in the form
"host:port". This will mirror the gopher content and optionally use HTML
templates:

-   `.template`: The template used for gopher menus
-   `.mdtemplate`: The template used for Markdown documents

