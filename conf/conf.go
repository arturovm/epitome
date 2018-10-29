package conf

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

var (
	Addr    string
	Port    int
	DataDir string
	Debug   bool
	Help    bool
)

func init() {
	flag.StringVarP(&Addr, "addr", "a", "127.0.0.1", "API server address")
	flag.IntVarP(&Port, "port", "p", 8080, "API server port")
	flag.StringVar(&DataDir, "data", "$HOME/.epitome", "Epitome data directory")
	flag.BoolVarP(&Debug, "debug", "d", false, "Enable debug mode")
	flag.BoolVarP(&Help, "help", "h", false, "Print this message")
	flag.Parse()
}

func PrintHelp() {
	fmt.Println("Usage: epitome [options]")
	flag.CommandLine.SortFlags = false
	flag.PrintDefaults()
}
