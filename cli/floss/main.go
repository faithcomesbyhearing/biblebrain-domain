package floss

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("abcdefg")
	auditCmd := flag.NewFlagSet("audit", flag.ExitOnError)
	auditBucket := auditCmd.String("bucket", "", "bucket name (eg dbp-prod)")
	auditFilesetId := auditCmd.String("filesetId", "", "filesetId to be audited")

	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeBucket := removeCmd.String("bucket", "", "bucket name (eg dbp-prod)")
	removeFilesetId := removeCmd.String("filesetId", "", "filesetId to be audited")

	if len(os.Args) < 2 {
		fmt.Println("expected 'audit' or 'enable' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "audit":
		auditCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'audit'")
		fmt.Println("  bucket:", *auditBucket)
		fmt.Println("  filesetId:", *auditFilesetId)
		fmt.Println("  tail:", auditCmd.Args())
		if len(*auditBucket) == 0 {
			*auditBucket = "dbp-staging"
		}
		if len(*auditFilesetId) == 0 {
			fmt.Println("provide filesetid as command line arg")
			return
		}
	case "remove":
		removeCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'remove'")
		fmt.Println("  bucket:", *removeBucket)
		fmt.Println("  filesetId:", *removeFilesetId)
		fmt.Println("  tail:", removeCmd.Args())
		if len(*removeBucket) == 0 {
			*removeBucket = "dbp-staging"
		}
		if len(*removeFilesetId) == 0 {
			fmt.Println("provide filesetid as command line arg")
			return
		}
	default:
		fmt.Println("expected 'audit' or 'remove' subcommands")
		os.Exit(1)
	}
}
