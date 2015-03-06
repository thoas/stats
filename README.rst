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

    import (
        "net/http"
        "github.com/thoas/stats"
    )

    func main() {
        h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte("{\"hello\": \"world\"}"))
        })

        handler = stats.New().Handler(h)
        http.ListenAndServe(":8080", handler)
    }

Usage
-----

Negroni
.......

If you are using negroni_ you can implement the handler as middleware:

.. code-block:: go

    package main

    import (
        "net/http"
        "github.com/codegangsta/negroni"
        "github.com/thoas/stats"
        "encoding/json"
    )

    func main() {
        middleware := stats.New()

        mux := http.NewServeMux()

        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte("{\"hello\": \"world\"}"))
        })

        mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")

            stats := middleware.GetStats()

            b, _ := json.Marshal(stats)

            w.Write(b)
        })

        n := negroni.Classic()
        n.Use(middleware)
        n.UseHandler(mux)
        n.Run(":3000")
    }

Inspiration
-----------

This reusable handler comes from a complete rip off of the great StatusMiddleware_
which comes from the awesome `go-json-rest`_.

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _StatusMiddleware: https://github.com/ant0ine/go-json-rest/blob/master/rest/status.go
.. _go-json-rest: https://github.com/ant0ine/go-json-rest
.. _negroni: https://github.com/codegangsta/negroni
