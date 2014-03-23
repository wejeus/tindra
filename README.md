
# Tindra - the glimmering site generator!

Like Rubys Jekyll but for Go! Project is in initial planning and hacking phase, stay tuned!

## Roadmap

* Add support for file inclusion
* Add template parsing
* Add export of generated site
* Add ability run as stand alone webserver
* Add filesytem watchers
* Add support for Disqus

## Documentation

Tindra relies on *convention over configuration* for many of its setttings.

### Special folders

	<site>/images 		- site specific media (assets used for building site framework)
	<site>/media		- content specific media (media used in posts)
	<site>/css

	<site>/includes		- snippets that can be included in other places (refer to an include by its filename). May contain subfolders.
	<site>/layouts
	<site>/posts

	404.html
	index.html
	config.yml

## API

Include files: {{.Include('filename')}}

## Post
First paragraph may contain a Yaml header. The header affects the rendering of the post. 

	---
	title: "My post title" // if not set uses filename as title
	layout: "layout.html"  // if not set just outputs parsed markdown
	---