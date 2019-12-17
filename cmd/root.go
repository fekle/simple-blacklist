package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fekle/simple-blacklist/pkg/blacklist"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var urlList []string
var outPath string
var exactPrefix string

// rootCmd defines the root command for this app
var rootCmd = &cobra.Command{
	Use:   "simple-blacklist",
	Short: "A simple tool to fetch, parse and merge domain blacklists",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// ensure output file
		outPath, err := filepath.Abs(outPath)
		if err != nil {
			log.Fatal(err)
		}
		outFile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		// init table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"source", "domains", "unique",
		})

		// process all sources
		var lists []*blacklist.Blacklist
		var listsMutex sync.Mutex
		var eg errgroup.Group
		for i := range urlList {
			blacklistURL := urlList[i]
			eg.Go(func() error {
				// init blacklist
				bl := blacklist.NewBlacklist(blacklistURL)

				// fetch and filter blacklist
				if err := bl.Process(exactPrefix); err != nil {
					return err
				}

				// add to list of blacklists
				listsMutex.Lock()
				lists = append(lists, bl)
				listsMutex.Unlock()

				// append info to table
				table.Append([]string{bl.Source, strconv.Itoa(bl.NTotal), strconv.Itoa(bl.NUniq)})

				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			log.Fatal(err)
		}

		// merge all blacklists into final list
		final := blacklist.Merge(lists)

		// show table with final list info
		table.Append([]string{""})
		table.Append([]string{outPath, strconv.Itoa(final.NTotal), strconv.Itoa(final.NUniq)})
		table.Render()

		// write final list to output file
		_, err = outFile.WriteString(strings.Join(final.Domains, "\n"))
		if err != nil {
			log.Fatal(err)
		}
	},
}

// set cobra options
func init() {
	rootCmd.Flags().StringSliceVarP(&urlList, "url", "u", []string{}, "comma-separated list of urls")
	rootCmd.MarkFlagRequired("url")

	rootCmd.Flags().StringVarP(&outPath, "output", "o", "", "path to write final blacklist to")
	rootCmd.MarkFlagRequired("outPath")

	rootCmd.Flags().StringVarP(&exactPrefix, "exactPrefix", "e", "", "prefix to prepend to non-wildcard entries")
}

// Execute runs the cobra root app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
