# Website Poller

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/SunSince90/website-poller/Go)
![GitHub top language](https://img.shields.io/github/languages/top/sunsince90/website-poller)
![GitHub](https://img.shields.io/github/license/sunsince90/website-poller)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sunsince90/website-poller)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/sunsince90/website-poller)

A simple and lightweight *Go* package that helps you make recurrent requests
to a website of your page.

## Overview

The package will poll a website every *X* seconds or at a random time at each
call with a provided list of *Headers* and a User Agent of your choice.
Alternatively, you can provide a list of user agents that can be rotated or
chosen randomly each time.
A *Handler Function* of your choice will be executed whenever the requests
completes - whether it failed or not.

### Features

* Load options from file or define them on your file
* Poll at a fixed time
* Poll at a random time based on a range of seconds to mimick user behavior,
i.e. between `[30 - 50]`
seconds
  * Example: first poll after 32 seconds
  * second poll after 45 seconds
  * third poll after 30 seconds
  * fourth poll after 37 seconds
  * and so on...
* Provide custom *Headers
* Provide a *User Agents* list with the ability to either:
  * Rotate them at each request
  * Pick a random one each time
* Provide no user agents list and let the package choose a random one each time

### Limitations and warnings

Remember that if you poll too aggressively you could be probably banned by the
website, put behind captchas or exceed quotas to the *API* service.
This package **will not** prevent you from being banned nor will solve captchas
for you. Remember to be polite and respect the rules defined by the website you
intend to poll.

The package does not support sending a body with each request yet.

### Features that will be introduced on future

* Headers generator to generate headers for every request
* Body generator to generate a different body for every request
* Custom http client
* Custom http request

## Install

```bash
go get github.com/SunSince90/website-poller
```

## How to use

First of all, import it in your go file:

```go
import (
    poller "github.com/SunSince90/website-poller"
)
```

Then, define a *Handler Function* that will be called when each request
completes.

```go
func handleResponse(id string, resp *http.Response, err error) {
    if err != nil {
        // handle the error here
    }

    // Do your stuff here...
}
```

Define the website to poll:

```go
// Poll a website every 30 seconds
page := &poller.Page {
    ID: "github-sunsince90",
    URL: "https://api.github.com/users/sunsince90",
}
```

Finally, start polling:

```go
p := poller.New(page)
p.SetHandlerFunc(handleResponse)

ctx, canc := context.WithCancel(context.Background())
p.Start(ctx)
```

## Examples

The above program will block the main thread, follow the examples contained
in the `examples` folder to learn more:

* [Log](./examples/log/log.go): a simple logger
* [File](./examples/file/file.go): load the pages to poll from a file
* [Custom](./examples/custom/custom.go): a more advanced poller with
polling options
* [Concurrent](./examples/concurrent/concurrent.go): how to load multiple
pollers and correctly wait for them to finish
