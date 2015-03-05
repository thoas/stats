Go stats handler
================

stats is a ``net/http`` handler in golang reporting various metrics about
your web application.

Installation
------------

1. Make sure you have a Go language compiler >= 1.3 (required) and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Download it:

::

    go get github.com/thoas/stats

5. Run your server:

.. code-block:: go

    package main

    func main() {
        // :)
    }

Usage
-----

...

Inspiration
-----------

This reusable handler comes from a complete rip off of the great StatusMiddleware_
which comes from the awesome `go-json-rest`_.

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _StatusMiddleware: https://github.com/ant0ine/go-json-rest/blob/master/rest/status.go
.. _go-json-rest: https://github.com/ant0ine/go-json-rest
