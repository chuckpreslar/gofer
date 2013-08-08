#gofer

__go·fer__ /ˈgō-fər/ _(n)_ - An employee who runs errands in addition to performing regular duties.

[![Build Status](https://drone.io/github.com/chuckpreslar/gofer/status.png)](https://drone.io/github.com/chuckpreslar/gofer/latest)

## About

You see too many Go packages relying on rake.  Why should we force people to use something that's not native to the platform they're developing for?  __gofer__ was made for just this reason.

As long as you place your task files within a `tasks` directory, and import the __gofer__ package, users will have access to them from anywhere as long as the binary is installed along their `PATH` environment variable.

## Installation

Assuming `/usr/local/bin` is on your path:

    $ go get -u github.com/chuckpreslar/gofer
    $ cd $GOPATH/src/github.com/chuckpreslar/gofer
    $ go build ./bin/gofer.go && mv ./gofer /usr/local/bin/gofer

## Usage

The following is the most basic example:

```go
//  $GOPATH/
//    src/
//      your_package/
//        tasks/
//          task_one.go

package tasks

import (
  "github.com/chuckpreslar/gofer"
)

var TaskOne = gofer.Register(gofer.Task{
  Label:       "task_one",
  Description: "Performs a simple task.",
  Action: func() error {
    // Perform something when called.
  },
})
```

To access the task you created, simply execute the following command from your terminal:

    $ gofer task_one

Yes, yes.. of course you can namespace (or sectionalize) your tasks:

```go
var TaskOne = gofer.Register(gofer.Task{
  Section:     "my_package:sub_section",
  Label:       "task_one",
  Description: "Performs a simple task.",
  Action: func() error {
    // Perform something when called.
  },
})

// $ gofer my_package:sub_section:task_one
```

Don't worry about the chicken or the egg problem; if you somehow manage to create a nested task before its parent, __gofer__ simply carves out a blank space for you to come back and potentially fill in later.

You can also give your tasks dependencies that will be evaluated in the appropriate order:

```go
var TaskTwo = gofer.Register(gofer.Task{
  Section:     "my_package",
  Label:       "task_two",
  Dependencies: []string{"my_package:task_two"},
  Description: "Performs a simple task after my_package:task_one executes.",
  Action: func() error {
    // Perform something when called.
  },
})

// $ gofer my_package:task_two # executes my_package:task_one followed by my_package:task_two
```

More complex dependency mapping is also resolved as long as there are no cyclic dependencies.

## Note

This is simply a proposal of how I feel a Go task manage should work.  If you feel you have a feature that should be added or that code is sloppy in an area, please, open a pull request for discussion.

## Documentation

View godoc or visit [godoc.org](http://godoc.org/github.com/chuckpreslar/gofer).

    $ godoc gofer

## License

> The MIT License (MIT)

> Copyright (c) 2013 Chuck Preslar

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.
