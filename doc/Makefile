all: build

puml:
	plantuml *.puml -tlatex:nopreamble

pdf:
	pandoc comparison.md -f markdown -t latex -o comparison.pdf --filter pandoc-citeproc --bibliography=references.bib --csl=numeric.csl

build: pdf puml

clean:
	rm -f *.tex *.pdf *.png *.svg
