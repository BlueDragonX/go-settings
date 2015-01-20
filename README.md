Go Settings in YAML
===================
The go-settings library implements a settings library on top of YAML. It is
meant to provide ease of use without resorting to type conversions and custom
YAML parsing.

Usage
-----
The library provides two functions for reading settings from YAML: `Read` and
`Load`. The `Read` function reads data from an `io.Reader` instance while the
`Load` function will load data from the provided file path. 

Using `Read`:

    import (
        "fmt"
        "github.com/BlueDragonX/go-settings"
        "os"
    )

    var file *os.File
    var err error
    var s *settings.Settings

    if file, err = os.Open; err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if s, err = settings.Read(rdr); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("got settings from reader")

Using `Load`:

    import (
        "fmt"
        "github.com/BlueDragonX/go-settings"
        "os"
    )

    if s, err = settings.Load("settings.yml"); if err == nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("got settings from file")
    
There is also `LoadOrExit` which does not return an error. It will call `Load`
and if it fails will print the error to stderr and exit. For example:

    import (
        "fmt"
        "github.com/BlueDragonX/go-settings"
    )

    s = settings.LoadOrExit("settings.yml")

    fmt.Println("got settings from file")

There are a handful of methods on the settings object which can be used to
retrieve values. All methods take a key string as the first argument. The key
string itself is a dot (.) separated list of object names. The different
methods will handle conversion to different types. With a few exceptions they
are self explanatory. These methods are:

- `Raw`: Return the raw value as an interface{}.
- `Object`: Return a settings object.
- `ObjectArray`: Return an array of settings objects.
- `String`
- `StringArray`
- `Int`
- `IntArray`
- `Float`
- `FloatArray`
- `Bool`
- `BoolArray`

The get methods may return a predefined error value to indicate failure. These are:

- `KeyError`: The key was not found.
- `TypeError`: Conversion to the requested type failed.

License
-------
Copyright (c) 2014 Ryan Bourgeois. Licensed under BSD-Modified. See the LICENSE
file for a copy of the license.
