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


Usage
-----

Basic net/http
..............

To use this handler directly with ``net/http``, you need to call the
middleware with the handler itself:

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

        handler := stats.New().Handler(h)
        http.ListenAndServe(":8080", handler)
    }

Negroni
.......

If you are using negroni_ you can implement the handler as
a simple middleware in ``server.go``:

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

Run it in a shell:

::

    $ go run server.go

Then in another shell run:

::

    $ curl http://localhost:3000/stats | python -m "json.tool"

Except the following result:

.. code-block:: json

    {
        "total_response_time": "1.907382ms",
        "average_response_time": "86.699\u00b5s",
        "average_response_time_sec": 8.6699e-05,
        "count": 1,
        "pid": 99894,
        "status_code_count": {
            "200": 1
        },
        "time": "2015-03-06 17:23:27.000677896 +0100 CET",
        "total_count": 22,
        "total_response_time_sec": 0.0019073820000000002,
        "total_status_code_count": {
            "200": 22
        },
        "unixtime": 1425659007,
        "uptime": "4m14.502271612s",
        "uptime_sec": 254.502271612
    }



See `examples <https://github.com/thoas/stats/blob/master/examples>`_ to
test them.



Inspiration
-----------

This reusable handler comes from a complete rip off of the great StatusMiddleware_
which is located in the `go-json-rest`_ repository.

Thanks to `Antoine Imbert <https://github.com/ant0ine>`_ for his work.

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _StatusMiddleware: https://github.com/ant0ine/go-json-rest/blob/master/rest/status.go
.. _go-json-rest: https://github.com/ant0ine/go-json-rest
.. _negroni: https://github.com/codegangsta/negroni
