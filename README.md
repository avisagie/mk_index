mk_index
========

Makes a simple index.html for a directory. It works ok for shallow
trees without many files. Folders can be collapsed by clicking on them.

Build
-----

	go get github.com/gobuffalo/packr/...
	bin/packr install mk_index
	
serve
=====

Dead simple http server that serves up the current working
directory. You can set Cache-Control: max-age and the listen address.

Build
-----

	go install serve

