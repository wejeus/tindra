
# Tindra - the glimmering site generator!

Like Rubys Jekyll but for Go! Project is in initial planning and hacking phase, stay tuned! Try it out by running:
	
	go run tindra.go res/testsite.se/

Should output a generated post on stdout. So far includes/layouts/markdown/data is implemented.


### Special folders

All files and folders in site root are copied verbatim except for a few that are special and one that is allways ignored (build/, config.yml). Files in root are parsed.


	<site>/images 		- site specific media (assets used for building site framework)
	<site>/media		- content specific media (media used in posts)
	<site>/css

	<site>/includes		- snippets that can be included in other places (refer to an include by its filename). May contain subfolders.
	<site>/layouts
	<site>/posts

	404.html
	index.html
	config.yml

## Generation dependency graph

Includes are static and nothing can be inserted
Layouts can use other layouts.
Posts can use layouts

## API 

#### Pre-processing 

Include files: {% include sidebar.html %}
For layouts: {% content %}
In posts: accesss variables using golang templates

## Post

All blog post files must begin with YAML front-matter. That is the first paragraph and contains a Yaml header. The header affects the rendering of the post. 

	---
	title: "My post title" // if not set uses filename as title
	layout: "layout.html"  // if not set just outputs parsed markdown
	---

