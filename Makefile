dl-htmx:
	curl https://unpkg.com/htmx.org@2.0.2/dist/htmx.min.js --create-dirs -o assets/js/htmx.min.js

dl-hyperscript:
	curl https://unpkg.com/hyperscript.org@0.9.12/dist/_hyperscript.min.js --create-dirs -o assets/js/_hyperscript.min.js

dl-assets: dl-htmx dl-hyperscript
