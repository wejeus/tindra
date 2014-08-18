
## Need support for (parsing of headers, should be put in post metadata)

date: "2013-05-08 23:46:11 +0200"
author: parkr 	/ 	author: all
version: 1.0.1
categories: [release]
per post permalink -> permalink: /my/custom/permalink
define default permalink using config.yml -> permalink: /news/:year/:month/:day/:title/

## Plugins (apply plugin by adding header field 'plugin: my_git_issue_plug' to handle thins like..)

* Add newer `language-` class name prefix to code blocks ([#1037][])

## Roadmap

* Add support for file inclusion
* Add ability run as stand alone webserver
* Add filesytem watchers
* Add support for Disqus
* Generation of feed.xml

## Refactor

* Namechange package context -> package site, type Site -> type Context
* It is the error implementation's responsibility to summarize the context. The error returned by os.Open formats as "open /etc/passwd: permission denied," not just "permission denied." The error returned by our Sqrt is missing information about the invalid argument.

## Documentation

Tindra relies on *convention over configuration* for many of its setttings.

## Notes

* Layouts are always keept preprocessed in memory and posts are stored with a pointer to layout and rendered when needed. Assumption here is that there will be more posts than layouts.

* When defining .yaml configfiles indendation must be done with spaces

* Halts directly on error

* http://www.elijahmanor.com/css-animated-hamburger-icon/

* func foo(args ...string) // args is treated as []string