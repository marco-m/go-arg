package arg

import (
	"fmt"
	"os"
	"strings"
)

func split(s string) []string {
	return strings.Split(s, " ")
}

// This example demonstrates basic usage
func Example() {
	// These are the args you would pass in on the command line
	os.Args = split("./example --foo=hello --bar")

	var args struct {
		Foo string
		Bar bool
	}
	MustParse(&args)
	fmt.Println(args.Foo, args.Bar)
	// output: hello true
}

// This example demonstrates arguments that have default values
func Example_defaultValues() {
	// These are the args you would pass in on the command line
	os.Args = split("./example")

	var args struct {
		Foo string
	}
	args.Foo = "default value"
	MustParse(&args)
	fmt.Println(args.Foo)
	// output: default value
}

// This example demonstrates arguments that are required
func Example_requiredArguments() {
	// These are the args you would pass in on the command line
	os.Args = split("./example --foo=abc --bar")

	var args struct {
		Foo string `arg:"required"`
		Bar bool
	}
	MustParse(&args)
	fmt.Println(args.Foo, args.Bar)
	// output: abc true
}

// This example demonstrates positional arguments
func Example_positionalArguments() {
	// These are the args you would pass in on the command line
	os.Args = split("./example in out1 out2 out3")

	var args struct {
		Input  string   `arg:"positional"`
		Output []string `arg:"positional"`
	}
	MustParse(&args)
	fmt.Println("In:", args.Input)
	fmt.Println("Out:", args.Output)
	// output:
	// In: in
	// Out: [out1 out2 out3]
}

// This example demonstrates arguments that have multiple values
func Example_multipleValues() {
	// The args you would pass in on the command line
	os.Args = split("./example --database localhost --ids 1 2 3")

	var args struct {
		Database string
		IDs      []int64
	}
	MustParse(&args)
	fmt.Printf("Fetching the following IDs from %s: %v", args.Database, args.IDs)
	// output: Fetching the following IDs from localhost: [1 2 3]
}

// This eample demonstrates multiple value arguments that can be mixed with
// other arguments.
func Example_multipleMixed() {
	os.Args = split("./example -c cmd1 db1 -f file1 db2 -c cmd2 -f file2 -f file3 db3 -c cmd3")
	var args struct {
		Commands  []string `arg:"-c,separate"`
		Files     []string `arg:"-f,separate"`
		Databases []string `arg:"positional"`
	}
	MustParse(&args)
	fmt.Println("Commands:", args.Commands)
	fmt.Println("Files:", args.Files)
	fmt.Println("Databases:", args.Databases)

	// output:
	// Commands: [cmd1 cmd2 cmd3]
	// Files: [file1 file2 file3]
	// Databases: [db1 db2 db3]
}

// This example shows the usage string generated by go-arg
func Example_helpText() {
	// These are the args you would pass in on the command line
	os.Args = split("./example --help")

	var args struct {
		Input    string   `arg:"positional"`
		Output   []string `arg:"positional"`
		Verbose  bool     `arg:"-v" help:"verbosity level"`
		Dataset  string   `help:"dataset to use"`
		Optimize int      `arg:"-O,help:optimization level"`
	}

	// This is only necessary when running inside golang's runnable example harness
	osExit = func(int) {}

	MustParse(&args)

	// output:
	// Usage: example [--verbose] [--dataset DATASET] [--optimize OPTIMIZE] INPUT [OUTPUT [OUTPUT ...]]
	//
	// Positional arguments:
	//   INPUT
	//   OUTPUT
	//
	// Options:
	//   --verbose, -v          verbosity level
	//   --dataset DATASET      dataset to use
	//   --optimize OPTIMIZE, -O OPTIMIZE
	//                          optimization level
	//   --help, -h             display this help and exit
}

// This example shows the usage string generated by go-arg when using subcommands
func Example_helpTextWithSubcommand() {
	// These are the args you would pass in on the command line
	os.Args = split("./example --help")

	type getCmd struct {
		Item string `arg:"positional" help:"item to fetch"`
	}

	type listCmd struct {
		Format string `help:"output format"`
		Limit  int
	}

	var args struct {
		Verbose bool
		Get     *getCmd  `arg:"subcommand" help:"fetch an item and print it"`
		List    *listCmd `arg:"subcommand" help:"list available items"`
	}

	// This is only necessary when running inside golang's runnable example harness
	osExit = func(int) {}

	MustParse(&args)

	// output:
	// Usage: example [--verbose]
	//
	// Options:
	//   --verbose
	//   --help, -h             display this help and exit
	//
	// Commands:
	//   get                    fetch an item and print it
	//   list                   list available items
}

// This example shows the usage string generated by go-arg when using subcommands
func Example_helpTextForSubcommand() {
	// These are the args you would pass in on the command line
	os.Args = split("./example get --help")

	type getCmd struct {
		Item string `arg:"positional" help:"item to fetch"`
	}

	type listCmd struct {
		Format string `help:"output format"`
		Limit  int
	}

	var args struct {
		Verbose bool
		Get     *getCmd  `arg:"subcommand" help:"fetch an item and print it"`
		List    *listCmd `arg:"subcommand" help:"list available items"`
	}

	// This is only necessary when running inside golang's runnable example harness
	osExit = func(int) {}

	MustParse(&args)

	// output:
	// Usage: example get ITEM
	//
	// Positional arguments:
	//   ITEM                   item to fetch
	//
	// Options:
	//   --help, -h             display this help and exit
}

// This example shows the error string generated by go-arg when an invalid option is provided
func Example_errorText() {
	// These are the args you would pass in on the command line
	os.Args = split("./example --optimize INVALID")

	var args struct {
		Input    string   `arg:"positional"`
		Output   []string `arg:"positional"`
		Verbose  bool     `arg:"-v" help:"verbosity level"`
		Dataset  string   `help:"dataset to use"`
		Optimize int      `arg:"-O,help:optimization level"`
	}

	// This is only necessary when running inside golang's runnable example harness
	osExit = func(int) {}
	stderr = os.Stdout

	MustParse(&args)

	// output:
	// Usage: example [--verbose] [--dataset DATASET] [--optimize OPTIMIZE] INPUT [OUTPUT [OUTPUT ...]]
	// error: error processing --optimize: strconv.ParseInt: parsing "INVALID": invalid syntax
}

// This example shows the error string generated by go-arg when an invalid option is provided
func Example_errorTextForSubcommand() {
	// These are the args you would pass in on the command line
	os.Args = split("./example get --count INVALID")

	type getCmd struct {
		Count int
	}

	var args struct {
		Get *getCmd `arg:"subcommand"`
	}

	// This is only necessary when running inside golang's runnable example harness
	osExit = func(int) {}
	stderr = os.Stdout

	MustParse(&args)

	// output:
	// Usage: example get [--count COUNT]
	// error: error processing --count: strconv.ParseInt: parsing "INVALID": invalid syntax
}

// This example demonstrates use of subcommands
func Example_subcommand() {
	// These are the args you would pass in on the command line
	os.Args = split("./example commit -a -m what-this-commit-is-about")

	type CheckoutCmd struct {
		Branch string `arg:"positional"`
		Track  bool   `arg:"-t"`
	}
	type CommitCmd struct {
		All     bool   `arg:"-a"`
		Message string `arg:"-m"`
	}
	type PushCmd struct {
		Remote      string `arg:"positional"`
		Branch      string `arg:"positional"`
		SetUpstream bool   `arg:"-u"`
	}
	var args struct {
		Checkout *CheckoutCmd `arg:"subcommand:checkout"`
		Commit   *CommitCmd   `arg:"subcommand:commit"`
		Push     *PushCmd     `arg:"subcommand:push"`
		Quiet    bool         `arg:"-q"` // this flag is global to all subcommands
	}

	// This is only necessary when running inside golang's runnable example harness
	osExit = func(int) {}
	stderr = os.Stdout

	MustParse(&args)

	switch {
	case args.Checkout != nil:
		fmt.Printf("checkout requested for branch %s\n", args.Checkout.Branch)
	case args.Commit != nil:
		fmt.Printf("commit requested with message \"%s\"\n", args.Commit.Message)
	case args.Push != nil:
		fmt.Printf("push requested from %s to %s\n", args.Push.Branch, args.Push.Remote)
	}

	// output:
	// commit requested with message "what-this-commit-is-about"
}
