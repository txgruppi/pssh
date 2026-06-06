package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func run() error {
	app := &cli.Command{
		Name:      "pssh",
		Usage:     "Run commands in parallel on multiple hosts via SSH",
		ArgsUsage: "<command|script file>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "hosts",
				Required: true,
				Usage:    "Path to JSON file with hosts definitions",
			},
			&cli.StringSliceFlag{
				Name:  "tag",
				Usage: "Filter hosts by tag (can be specified multiple times)",
			},
			&cli.BoolFlag{
				Name:  "sudo",
				Usage: "Run the command with sudo",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Do not do anything, just print what would be done",
			},
			&cli.IntFlag{
				Name:  "limit",
				Value: -1,
				Usage: "Limit the number of concurrent connections (default: unlimited)",
			},
			&cli.BoolFlag{
				Name:  "no-known-hosts",
				Usage: "Do not check known hosts",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			hosts, err := loadHostsFromFile(c.String("hosts"))
			if err != nil {
				return err
			}
			hosts = hosts.FilterByTags(c.StringSlice("tag"))
			isFile := c.Args().Len() == 1 && isFile(c.Args().Get(0))
			var filepath string
			if isFile {
				filepath, err = copyTempFile(c.Args().Get(0))
				if err != nil {
					return err
				}
			} else {
				filepath, err = makeTempExecFile([]byte(strings.Join(c.Args().Slice(), " ")))
				if err != nil {
					return err
				}
			}
			defer os.Remove(filepath)
			var m sync.Mutex
			g, _ := errgroup.WithContext(ctx)
			g.SetLimit(c.Int("limit"))
			for i, host := range hosts {
				g.Go(func() error {
					stdout, stderr, err := sshUploadRunAndCleanup(host, filepath, c.Bool("sudo"), c.Bool("dry-run"), c.Bool("no-known-hosts"))
					m.Lock()
					defer m.Unlock()
					if len(stdout) > 0 {
						stdout = prefixLines([]byte(hosts[i].Name+": "), stdout)
					}
					if len(stderr) > 0 {
						stderr = prefixLines([]byte(hosts[i].Name+": "), stderr)
					}
					if len(stdout) > 0 {
						os.Stdout.Write(stdout)
					}
					if len(stderr) > 0 {
						os.Stderr.Write(stderr)
					}
					if err != nil {
						return err
					}
					return nil
				})
			}
			err = g.Wait()
			return err
		},
	}
	return app.Run(context.Background(), os.Args)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
