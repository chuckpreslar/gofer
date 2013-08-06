#gofer

__go·fer__ /ˈgō-fər/ _(n)_ - An employee who runs errands in addition to performing regular duties.

[![Build Status](https://drone.io/github.com/chuckpreslar/gofer/status.png)](https://drone.io/github.com/chuckpreslar/gofer/latest)

## About

__gofer__ was developed from a need to assist user's in preforming initial setup tasks when starting with a new project in Go. Though not as extensive as tools like rake for Ruby or grunt for Javascript (yet), __gofer__ is still plenty useful.

#### For Package Maintainers
-------------------------

If you maintain a package, to make use of gofer you simply need to place your "tasks" within a directory named `tasks` somewhere inside your projects source directory.  When a user installs your package (having __gofer__ installed), all they'll need to do in order to access your task is call for it.

    $ gofer <namespace>:<task>


## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get -u github.com/chuckpreslar/gofer

