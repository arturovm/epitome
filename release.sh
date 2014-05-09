#!/bin/sh


ProgName=$(basename $0)
sub_help(){
	echo "Usage: $ProgName <subcommand> [options]\n"
	echo "Subcommands:"
	echo " mac windows linux"
	echo ""
	echo "For help with each subcommand run:"
	echo "$ProgName <subcommand> -h|--help"
	echo ""
}

build() {
	echo "Building $GOOS-$GOARCH..."
	git checkout master
	go build -v -o pond *.go
	mkdir release
	cd release
	mkdir pond
	mv ../pond pond/
	cp -r ../static pond/
	zip -r pond_${GOOS}_${GOARCH}.zip pond
}

cleanup() {
	echo "Cleaning up..."
	rm -rf pond
}

sub_linux(){
	export GOARCH=amd64
	export GOOS=linux
	build
	cleanup
	echo "Done"
}

sub_mac(){
	export GOARCH=amd64
	export GOOS=darwin
	build
	cleanup
	echo "Done."
}

sub_windows(){
	export GOARCH=amd64
	export GOOS=windows
	build
	cleanup
	echo "Done."
}

subcommand=$1
case $subcommand in
	"" | "-h" | "--help")
		sub_help
		;;
	*)
		shift
		sub_${subcommand} $@
		if [ $? = 127 ]; then
			echo "Error: '$subcommand' is not a known subcommand." >&2
			echo " Run '$ProgName --help' for a list of known subcommands." >&2
			exit 1
		fi
		;;
esac
