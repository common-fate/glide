// +build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	
)

func main() {
	// Use local types and functions in order to avoid name conflicts with additional magefiles.
	type arguments struct {
		Verbose       bool          // print out log statements
		List          bool          // print out a list of targets
		Help          bool          // print out help for a specific target
		Timeout       time.Duration // set a timeout to running the targets
		Args          []string      // args contain the non-flag command-line arguments
	}

	parseBool := func(env string) bool {
		val := os.Getenv(env)
		if val == "" {
			return false
		}		
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.Printf("warning: environment variable %s is not a valid bool value: %v", env, val)
			return false
		}
		return b
	}

	parseDuration := func(env string) time.Duration {
		val := os.Getenv(env)
		if val == "" {
			return 0
		}		
		d, err := time.ParseDuration(val)
		if err != nil {
			log.Printf("warning: environment variable %s is not a valid duration value: %v", env, val)
			return 0
		}
		return d
	}
	args := arguments{}
	fs := flag.FlagSet{}
	fs.SetOutput(os.Stdout)

	// default flag set with ExitOnError and auto generated PrintDefaults should be sufficient
	fs.BoolVar(&args.Verbose, "v", parseBool("MAGEFILE_VERBOSE"), "show verbose output when running targets")
	fs.BoolVar(&args.List, "l", parseBool("MAGEFILE_LIST"), "list targets for this binary")
	fs.BoolVar(&args.Help, "h", parseBool("MAGEFILE_HELP"), "print out help for a specific target")
	fs.DurationVar(&args.Timeout, "t", parseDuration("MAGEFILE_TIMEOUT"), "timeout in duration parsable format (e.g. 5m30s)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stdout, `
%s [options] [target]

Commands:
  -l    list targets in this binary
  -h    show this help

Options:
  -h    show description of a target
  -t <string>
        timeout in duration parsable format (e.g. 5m30s)
  -v    show verbose output when running targets
 `[1:], filepath.Base(os.Args[0]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag will have printed out an error already.
		return
	}
	args.Args = fs.Args()
	if args.Help && len(args.Args) == 0 {
		fs.Usage()
		return
	}
		
	// color is ANSI color type
	type color int

	// If you add/change/remove any items in this constant,
	// you will need to run "stringer -type=color" in this directory again.
	// NOTE: Please keep the list in an alphabetical order.
	const (
		black color = iota
		red
		green
		yellow
		blue
		magenta
		cyan
		white
		brightblack
		brightred
		brightgreen
		brightyellow
		brightblue
		brightmagenta
		brightcyan
		brightwhite
	)

	// AnsiColor are ANSI color codes for supported terminal colors.
	var ansiColor = map[color]string{
		black:         "\u001b[30m",
		red:           "\u001b[31m",
		green:         "\u001b[32m",
		yellow:        "\u001b[33m",
		blue:          "\u001b[34m",
		magenta:       "\u001b[35m",
		cyan:          "\u001b[36m",
		white:         "\u001b[37m",
		brightblack:   "\u001b[30;1m",
		brightred:     "\u001b[31;1m",
		brightgreen:   "\u001b[32;1m",
		brightyellow:  "\u001b[33;1m",
		brightblue:    "\u001b[34;1m",
		brightmagenta: "\u001b[35;1m",
		brightcyan:    "\u001b[36;1m",
		brightwhite:   "\u001b[37;1m",
	}
	
	const _color_name = "blackredgreenyellowbluemagentacyanwhitebrightblackbrightredbrightgreenbrightyellowbrightbluebrightmagentabrightcyanbrightwhite"

	var _color_index = [...]uint8{0, 5, 8, 13, 19, 23, 30, 34, 39, 50, 59, 70, 82, 92, 105, 115, 126}

	colorToLowerString := func (i color) string {
		if i < 0 || i >= color(len(_color_index)-1) {
			return "color(" + strconv.FormatInt(int64(i), 10) + ")"
		}
		return _color_name[_color_index[i]:_color_index[i+1]]
	}

	// ansiColorReset is an ANSI color code to reset the terminal color.
	const ansiColorReset = "\033[0m"

	// defaultTargetAnsiColor is a default ANSI color for colorizing targets.
	// It is set to Cyan as an arbitrary color, because it has a neutral meaning
	var defaultTargetAnsiColor = ansiColor[cyan]

	getAnsiColor := func(color string) (string, bool) {
		colorLower := strings.ToLower(color)
		for k, v := range ansiColor {
			colorConstLower := colorToLowerString(k)
			if colorConstLower == colorLower {
				return v, true
			}
		}
		return "", false
	}

	// Terminals which  don't support color:
	// 	TERM=vt100
	// 	TERM=cygwin
	// 	TERM=xterm-mono
    var noColorTerms = map[string]bool{
		"vt100":      false,
		"cygwin":     false,
		"xterm-mono": false,
	}

	// terminalSupportsColor checks if the current console supports color output
	//
	// Supported:
	// 	linux, mac, or windows's ConEmu, Cmder, putty, git-bash.exe, pwsh.exe
	// Not supported:
	// 	windows cmd.exe, powerShell.exe
	terminalSupportsColor := func() bool {
		envTerm := os.Getenv("TERM")
		if _, ok := noColorTerms[envTerm]; ok {
			return false
		}
		return true
	}

	// enableColor reports whether the user has requested to enable a color output.
	enableColor := func() bool {
		b, _ := strconv.ParseBool(os.Getenv("MAGEFILE_ENABLE_COLOR"))
		return b
	}

	// targetColor returns the ANSI color which should be used to colorize targets.
	targetColor := func() string {
		s, exists := os.LookupEnv("MAGEFILE_TARGET_COLOR")
		if exists == true {
			if c, ok := getAnsiColor(s); ok == true {
				return c
			}
		}
		return defaultTargetAnsiColor
	}

	// store the color terminal variables, so that the detection isn't repeated for each target
	var enableColorValue = enableColor() && terminalSupportsColor()
	var targetColorValue = targetColor()

	printName := func(str string) string {
		if enableColorValue {
			return fmt.Sprintf("%s%s%s", targetColorValue, str, ansiColorReset)
		} else {
			return str
		}
	}

	list := func() error {
		
		targets := map[string]string{
			"build:accessHandler": "",
			"build:backend": "builds the Go API for the AWS Lambda runtime.",
			"build:eventHandler": "",
			"build:frontend": "generates the React static frontend.",
			"build:frontendAWSExports": "",
			"build:frontendDeployer": "",
			"build:granter": "",
			"build:slackNotifier": "",
			"build:syncer": "",
			"build:webhook": "",
			"clean": "removes build and packaging artifacts.",
			"deploy:cdk": "deploys the CDK infrastructure stack to AWS",
			"deploy:dev": "provisions a development environment",
			"deploy:frontend": "uploads the frontend to S3 and invalidates CloudFront",
			"deploy:production": "",
			"deploy:staging": "provisions a staging environment env should be 'dev' or 'test' to match a CDK internal deployment environment",
			"deploy:stagingCDK": "deploys a staging version of the CDK infrastructure.",
			"deploy:stagingDNS": "sets a DNS CNAME entry in Route53 pointing to the CloudFront domain.",
			"deploy:stagingFrontend": "uploads the frontend to the S3 bucket and invalidates CloudFront.",
			"deps:npm": "installs NPM dependencies for the repository using pnpm.",
			"destroy": "deprovisions the CDK stack.",
			"devConfig": "sets up the granted-deployment.yml file",
			"dotenv": "updates the .env file based on the deployed CDK infrastructure",
			"package": "",
			"packageAccessHandler": "zips the Go access handler API so that it can be deployed to Lambda.",
			"packageBackend": "zips the Go API so that it can be deployed to Lambda.",
			"packageEventHandler": "zips the Go event handler so that it can be deployed to Lambda.",
			"packageFrontendDeployer": "zips the Go frontend deployer so that it can be deployed to Lambda.",
			"packageGranter": "zips the Go granter so that it can be deployed to Lambda.",
			"packageSlackNotifier": "PackageNotifier zips the Go notifier so that it can be deployed to Lambda.",
			"packageSyncer": "zips the Go Syncer function handler so that it can be deployed to Lambda.",
			"packageWebhook": "zips the Go webhook handler so that it can be deployed to Lambda.",
			"release:production": "",
			"release:productionCDK": "",
			"release:publishCDKAssets": "",
			"release:publishCloudFormation": "",
			"release:publishFrontendAssets": "",
			"release:publishManifest": "updates the manifest.json file in the release bucket with the latest version information, so that our customer deployment tooling knows there is a new version available.",
			"watch": "hot-reloads the CDK deployment when local files change.",
		}

		keys := make([]string, 0, len(targets))
		for name := range targets {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		fmt.Println("Targets:")
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		for _, name := range keys {
			fmt.Fprintf(w, "  %v\t%v\n", printName(name), targets[name])
		}
		err := w.Flush()
		return err
	}

	var ctx context.Context
	var ctxCancel func()

	getContext := func() (context.Context, func()) {
		if ctx != nil {
			return ctx, ctxCancel
		}

		if args.Timeout != 0 {
			ctx, ctxCancel = context.WithTimeout(context.Background(), args.Timeout)
		} else {
			ctx = context.Background()
			ctxCancel = func() {}
		}
		return ctx, ctxCancel
	}

	runTarget := func(fn func(context.Context) error) interface{} {
		var err interface{}
		ctx, cancel := getContext()
		d := make(chan interface{})
		go func() {
			defer func() {
				err := recover()
				d <- err
			}()
			err := fn(ctx)
			d <- err
		}()
		select {
		case <-ctx.Done():
			cancel()
			e := ctx.Err()
			fmt.Printf("ctx err: %v\n", e)
			return e
		case err = <-d:
			cancel()
			return err
		}
	}
	// This is necessary in case there aren't any targets, to avoid an unused
	// variable error.
	_ = runTarget

	handleError := func(logger *log.Logger, err interface{}) {
		if err != nil {
			logger.Printf("Error: %+v\n", err)
			type code interface {
				ExitStatus() int
			}
			if c, ok := err.(code); ok {
				os.Exit(c.ExitStatus())
			}
			os.Exit(1)
		}
	}
	_ = handleError

	// Set MAGEFILE_VERBOSE so mg.Verbose() reflects the flag value.
	if args.Verbose {
		os.Setenv("MAGEFILE_VERBOSE", "1")
	} else {
		os.Setenv("MAGEFILE_VERBOSE", "0")
	}

	log.SetFlags(0)
	if !args.Verbose {
		log.SetOutput(ioutil.Discard)
	}
	logger := log.New(os.Stderr, "", 0)
	if args.List {
		if err := list(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	if args.Help {
		if len(args.Args) < 1 {
			logger.Println("no target specified")
			os.Exit(2)
		}
		switch strings.ToLower(args.Args[0]) {
			case "build:accesshandler":
				
				fmt.Print("Usage:\n\n\tmage build:accesshandler\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:backend":
				fmt.Println("Backend builds the Go API for the AWS Lambda runtime.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage build:backend\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:eventhandler":
				
				fmt.Print("Usage:\n\n\tmage build:eventhandler\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:frontend":
				fmt.Println("Frontend generates the React static frontend.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage build:frontend\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:frontendawsexports":
				
				fmt.Print("Usage:\n\n\tmage build:frontendawsexports\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:frontenddeployer":
				
				fmt.Print("Usage:\n\n\tmage build:frontenddeployer\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:granter":
				
				fmt.Print("Usage:\n\n\tmage build:granter\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:slacknotifier":
				
				fmt.Print("Usage:\n\n\tmage build:slacknotifier\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:syncer":
				
				fmt.Print("Usage:\n\n\tmage build:syncer\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "build:webhook":
				
				fmt.Print("Usage:\n\n\tmage build:webhook\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "clean":
				fmt.Println("Clean removes build and packaging artifacts.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage clean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:cdk":
				fmt.Println("CDK deploys the CDK infrastructure stack to AWS")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:cdk\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:dev":
				fmt.Println("Dev provisions a development environment")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:dev\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:frontend":
				fmt.Println("Frontend uploads the frontend to S3 and invalidates CloudFront")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:frontend\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:production":
				
				fmt.Print("Usage:\n\n\tmage deploy:production <releaseBucket> <versionHash> <stackName> <cfnParamsJSON>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:staging":
				fmt.Println("Staging provisions a staging environment env should be 'dev' or 'test' to match a CDK internal deployment environment")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:staging <env> <name>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:stagingcdk":
				fmt.Println("StagingCDK deploys a staging version of the CDK infrastructure.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:stagingcdk <env> <name>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:stagingdns":
				fmt.Println("StagingDNS sets a DNS CNAME entry in Route53 pointing to the CloudFront domain.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:stagingdns <env> <name>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deploy:stagingfrontend":
				fmt.Println("StagingFrontend uploads the frontend to the S3 bucket and invalidates CloudFront. It requires an internal deployment environment and name to be specified.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deploy:stagingfrontend <env> <name>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "deps:npm":
				fmt.Println("NPM installs NPM dependencies for the repository using pnpm.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage deps:npm\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "destroy":
				fmt.Println("Destroy deprovisions the CDK stack.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage destroy\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "devconfig":
				fmt.Println("DevConfig sets up the granted-deployment.yml file")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage devconfig\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "dotenv":
				fmt.Println("Dotenv updates the .env file based on the deployed CDK infrastructure")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage dotenv\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "package":
				
				fmt.Print("Usage:\n\n\tmage package\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packageaccesshandler":
				fmt.Println("PackageAccessHandler zips the Go access handler API so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packageaccesshandler\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packagebackend":
				fmt.Println("PackageBackend zips the Go API so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packagebackend\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packageeventhandler":
				fmt.Println("PackageEventHandler zips the Go event handler so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packageeventhandler\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packagefrontenddeployer":
				fmt.Println("PackageFrontendDeployer zips the Go frontend deployer so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packagefrontenddeployer\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packagegranter":
				fmt.Println("PackageGranter zips the Go granter so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packagegranter\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packageslacknotifier":
				fmt.Println("PackageNotifier zips the Go notifier so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packageslacknotifier\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packagesyncer":
				fmt.Println("PackageSyncer zips the Go Syncer function handler so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packagesyncer\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "packagewebhook":
				fmt.Println("PackageWebhook zips the Go webhook handler so that it can be deployed to Lambda.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage packagewebhook\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "release:production":
				
				fmt.Print("Usage:\n\n\tmage release:production <releaseBucket> <versionHash>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "release:productioncdk":
				
				fmt.Print("Usage:\n\n\tmage release:productioncdk <releaseBucket> <versionHash>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "release:publishcdkassets":
				
				fmt.Print("Usage:\n\n\tmage release:publishcdkassets <releaseBucket> <versionHash>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "release:publishcloudformation":
				
				fmt.Print("Usage:\n\n\tmage release:publishcloudformation <releaseBucket> <versionHash>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "release:publishfrontendassets":
				
				fmt.Print("Usage:\n\n\tmage release:publishfrontendassets <releaseBucket> <versionHash>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "release:publishmanifest":
				fmt.Println("PublishManifest updates the manifest.json file in the release bucket with the latest version information, so that our customer deployment tooling knows there is a new version available.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage release:publishmanifest <releaseBucket> <version>\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			case "watch":
				fmt.Println("Watch hot-reloads the CDK deployment when local files change.")
				fmt.Println()
				
				fmt.Print("Usage:\n\n\tmage watch\n\n")
				var aliases []string
				if len(aliases) > 0 {
					fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
				}
				return
			default:
				logger.Printf("Unknown target: %q\n", args.Args[0])
				os.Exit(2)
		}
	}
	if len(args.Args) < 1 {
		if err := list(); err != nil {
			logger.Println("Error:", err)
			os.Exit(1)
		}
		return
	}
	for x := 0; x < len(args.Args); {
		target := args.Args[x]
		x++

		// resolve aliases
		switch strings.ToLower(target) {
		
		}

		switch strings.ToLower(target) {
		
			case "build:accesshandler":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:AccessHandler\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:AccessHandler")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.AccessHandler()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:backend":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:Backend\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:Backend")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.Backend()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:eventhandler":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:EventHandler\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:EventHandler")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.EventHandler()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:frontend":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:Frontend\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:Frontend")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.Frontend()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:frontendawsexports":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:FrontendAWSExports\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:FrontendAWSExports")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.FrontendAWSExports()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:frontenddeployer":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:FrontendDeployer\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:FrontendDeployer")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.FrontendDeployer()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:granter":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:Granter\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:Granter")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.Granter()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:slacknotifier":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:SlackNotifier\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:SlackNotifier")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.SlackNotifier()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:syncer":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:Syncer\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:Syncer")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.Syncer()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "build:webhook":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:Webhook\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:Webhook")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.Webhook()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "clean":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Clean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Clean")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Clean()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:cdk":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:CDK\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:CDK")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.CDK()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:dev":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:Dev\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:Dev")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.Dev()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:frontend":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:Frontend\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:Frontend")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.Frontend()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:production":
				expected := x + 4
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:Production\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:Production")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
			arg2 := args.Args[x]
			x++
			arg3 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.Production(ctx, arg0, arg1, arg2, arg3)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:staging":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:Staging\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:Staging")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					Deploy{}.Staging(arg0, arg1)
					return nil
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:stagingcdk":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:StagingCDK\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:StagingCDK")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.StagingCDK(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:stagingdns":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:StagingDNS\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:StagingDNS")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.StagingDNS(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deploy:stagingfrontend":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deploy:StagingFrontend\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deploy:StagingFrontend")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Deploy{}.StagingFrontend(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "deps:npm":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Deps:NPM\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Deps:NPM")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Deps{}.NPM()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "destroy":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Destroy\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Destroy")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Destroy()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "devconfig":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"DevConfig\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "DevConfig")
				}
				
				wrapFn := func(ctx context.Context) error {
					return DevConfig()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "dotenv":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Dotenv\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Dotenv")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Dotenv()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "package":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Package\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Package")
				}
				
				wrapFn := func(ctx context.Context) error {
					Package()
					return nil
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packageaccesshandler":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageAccessHandler\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageAccessHandler")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageAccessHandler()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packagebackend":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageBackend\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageBackend")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageBackend()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packageeventhandler":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageEventHandler\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageEventHandler")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageEventHandler()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packagefrontenddeployer":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageFrontendDeployer\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageFrontendDeployer")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageFrontendDeployer()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packagegranter":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageGranter\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageGranter")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageGranter()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packageslacknotifier":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageSlackNotifier\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageSlackNotifier")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageSlackNotifier()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packagesyncer":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageSyncer\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageSyncer")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageSyncer()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "packagewebhook":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"PackageWebhook\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "PackageWebhook")
				}
				
				wrapFn := func(ctx context.Context) error {
					return PackageWebhook()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "release:production":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Release:Production\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Release:Production")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					Release{}.Production(arg0, arg1)
					return nil
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "release:productioncdk":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Release:ProductionCDK\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Release:ProductionCDK")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Release{}.ProductionCDK(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "release:publishcdkassets":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Release:PublishCDKAssets\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Release:PublishCDKAssets")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Release{}.PublishCDKAssets(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "release:publishcloudformation":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Release:PublishCloudFormation\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Release:PublishCloudFormation")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Release{}.PublishCloudFormation(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "release:publishfrontendassets":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Release:PublishFrontendAssets\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Release:PublishFrontendAssets")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Release{}.PublishFrontendAssets(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "release:publishmanifest":
				expected := x + 2
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Release:PublishManifest\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Release:PublishManifest")
				}
				
			arg0 := args.Args[x]
			x++
			arg1 := args.Args[x]
			x++
				wrapFn := func(ctx context.Context) error {
					return Release{}.PublishManifest(arg0, arg1)
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
			case "watch":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Watch\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Watch")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Watch()
				}
				ret := runTarget(wrapFn)
				handleError(logger, ret)
		
		default:
			logger.Printf("Unknown target specified: %q\n", target)
			os.Exit(2)
		}
	}
}




