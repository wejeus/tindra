## Roadmap

* Add support for file inclusion
* Add template parsing
* Add export of generated site
* Add ability run as stand alone webserver
* Add filesytem watchers
* Add support for Disqus

## Documentation

Tindra relies on *convention over configuration* for many of its setttings.

## Notes

* Layouts are always keept preprocessed in memory and posts are stored with a pointer to layout and rendered when needed. Assumption here is that there will be more posts than layouts.

* When defining .yaml configfiles indendation must be done with spaces

* Halts directly on error

* http://www.elijahmanor.com/css-animated-hamburger-icon/