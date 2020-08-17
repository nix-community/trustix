all: build

puml:
	plantuml *.puml -tlatex:nopreamble

build: puml

clean:
	rm -f *.tex *.pdf *.png *.svg
