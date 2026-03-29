package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/Joeri-Abbo/taco-lib/vulndb"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var (
		sourcesFlag string
		fullFlag    bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the vulnerability database from multiple sources",
		Long: `Update the vulnerability database by fetching from multiple sources.

By default, all sources are enabled: nvd, osv, ghsa, alpine-secdb, debian,
ubuntu, redhat, alas, cisa-kev.

On first run (no existing cache), a full historical fetch is performed
automatically. Subsequent runs fetch only the last 7 days incrementally.

Use --full to force a complete re-fetch of all historical data.

Use --sources to select specific sources:
  taco-vulndb update --sources nvd,osv,ghsa
  taco-vulndb update --full

Environment variables for API keys:
  TACO_NVD_API_KEY  — NVD API key (increases rate limit)
  GITHUB_TOKEN      — GitHub token for GHSA (increases rate limit)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			var sources []vulndb.SourceFetcher
			if sourcesFlag != "" {
				var names []vulndb.SourceName
				for _, s := range strings.Split(sourcesFlag, ",") {
					names = append(names, vulndb.SourceName(strings.TrimSpace(s)))
				}
				sources = vulndb.NewFetchersForSources(names)
			}

			if !quiet {
				mode := "incremental"
				if fullFlag || !cache.Exists() {
					mode = "full"
				}
				if len(sources) > 0 {
					names := make([]string, len(sources))
					for i, s := range sources {
						names[i] = string(s.Name())
					}
					slog.Info("updating vulnerability database...", "mode", mode, "sources", strings.Join(names, ", "))
				} else {
					slog.Info("updating vulnerability database from all sources...", "mode", mode)
				}
			}

			progressFn := func(source string, fetched, total int) {
				if !quiet {
					if total > 0 {
						fmt.Fprintf(os.Stderr, "\r[%s] Fetching: %d / %d", source, fetched, total)
					} else {
						fmt.Fprintf(os.Stderr, "\r[%s] Fetched: %d entries", source, fetched)
					}
				}
			}

			opts := &vulndb.UpdateOptions{Full: fullFlag}
			if err := vulndb.UpdateMultiSource(ctx, cache, sources, opts, progressFn); err != nil {
				return fmt.Errorf("updating database: %w", err)
			}

			if !quiet {
				fmt.Fprintln(os.Stderr)
				slog.Info("database updated successfully", "path", cache.DBPath())
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&sourcesFlag, "sources", "", "comma-separated list of sources (default: all)")
	cmd.Flags().BoolVar(&fullFlag, "full", false, "force full historical fetch instead of incremental")

	return cmd
}

func newDownloadCmd() *cobra.Command {
	var url string

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a pre-built database from a URL",
		Long: `Download a pre-built vulnerability database from a URL.

Examples:
  taco-vulndb download --url https://example.com/vulndb.json
  taco-vulndb download --url https://example.com/vulndb.json.gz`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				return fmt.Errorf("--url is required")
			}

			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			if !quiet {
				slog.Info("downloading vulnerability database...", "url", url)
			}

			if err := vulndb.DownloadDB(ctx, cache, url, nil); err != nil {
				return fmt.Errorf("downloading database: %w", err)
			}

			meta, _ := cache.ReadMeta()
			if !quiet {
				slog.Info("database downloaded successfully",
					"path", cache.DBPath(),
					"entries", meta.EntryCount,
				)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "URL to download the database from (required)")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}

func newLoadCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Import a local database file into the cache",
		Long: `Load a vulnerability database from a local file into the cache.

Examples:
  taco-vulndb load --file /path/to/vulndb.json
  taco-vulndb load --file /path/to/vulndb.json.gz`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file is required")
			}

			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			if !quiet {
				slog.Info("loading vulnerability database...", "file", file)
			}

			if err := vulndb.LoadDBFromFile(cache, file); err != nil {
				return fmt.Errorf("loading database: %w", err)
			}

			meta, _ := cache.ReadMeta()
			if !quiet {
				slog.Info("database loaded successfully",
					"path", cache.DBPath(),
					"entries", meta.EntryCount,
				)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "path to database file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newBuildCmd() *cobra.Command {
	var (
		output string
		days   int
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Fetch from NVD and build a portable database file",
		Long: `Build a standalone vulnerability database file by fetching from NVD.

Examples:
  taco-vulndb build --output ./vulndb/vulndb.json
  taco-vulndb build --output ./vulndb/vulndb.json --days 30`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if output == "" {
				return fmt.Errorf("--output is required")
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			if !quiet {
				slog.Info("building vulnerability database...", "output", output, "days", days)
			}

			progressFn := func(fetched, total int) {
				if !quiet {
					fmt.Fprintf(os.Stderr, "\rFetching vulnerabilities: %d / %d", fetched, total)
				}
			}

			if err := vulndb.BuildDB(ctx, output, days, progressFn); err != nil {
				return fmt.Errorf("building database: %w", err)
			}

			if !quiet {
				fmt.Fprintln(os.Stderr)
				slog.Info("database built successfully", "output", output)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	cmd.Flags().IntVar(&days, "days", 120, "number of days of CVE history to fetch")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func newExportCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the cached database as a gzip file for distribution",
		Long: `Export the local cached database as a compressed file.

Examples:
  taco-vulndb export --output vulndb.json.gz`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if output == "" {
				return fmt.Errorf("--output is required")
			}

			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			if err := vulndb.ExportGzip(cache, output); err != nil {
				return fmt.Errorf("exporting database: %w", err)
			}

			info, _ := os.Stat(output)
			if !quiet && info != nil {
				slog.Info("database exported",
					"output", output,
					"size_mb", fmt.Sprintf("%.1f", float64(info.Size())/(1024*1024)),
				)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func newServeCmd() *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start an HTTP server to host the database",
		Long: `Start a lightweight HTTP server that hosts the cached vulnerability database.

Endpoints:
  GET /vulndb.json     — database file (JSON)
  GET /vulndb.json.gz  — database file (gzip-compressed)
  GET /meta.json       — database metadata
  GET /health          — health check

Examples:
  taco-vulndb serve
  taco-vulndb serve --addr :9090`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			return vulndb.Serve(ctx, vulndb.ServeOptions{
				Addr:  addr,
				Cache: cache,
			})
		},
	}

	cmd.Flags().StringVar(&addr, "addr", ":8080", "listen address")

	return cmd
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show vulnerability database status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			if !cache.Exists() {
				fmt.Println("No vulnerability database found.")
				fmt.Println()
				fmt.Println("Populate it using one of:")
				fmt.Println("  taco-vulndb update                  — fetch from vulnerability sources")
				fmt.Println("  taco-vulndb download --url <url>    — download pre-built DB")
				fmt.Println("  taco-vulndb load --file <path>      — import a local DB file")
				return nil
			}

			meta, err := cache.ReadMeta()
			if err != nil {
				return fmt.Errorf("reading cache metadata: %w", err)
			}

			stale, _ := cache.IsStale()

			fmt.Printf("Database path:    %s\n", cache.DBPath())
			fmt.Printf("Last updated:     %s\n", meta.LastUpdated.Format("2006-01-02 15:04:05 MST"))
			fmt.Printf("Entry count:      %d\n", meta.EntryCount)
			if meta.SourceURL != "" {
				fmt.Printf("Source:           %s\n", meta.SourceURL)
			}
			if meta.ETag != "" {
				fmt.Printf("ETag:             %s\n", meta.ETag)
			}
			if stale {
				fmt.Printf("Status:           STALE (older than %s)\n", cache.MaxAge)
			} else {
				fmt.Printf("Status:           OK\n")
			}

			if len(meta.Sources) > 0 {
				fmt.Println()
				fmt.Println("Sources:")
				for name, sm := range meta.Sources {
					fmt.Printf("  %-15s  %d entries  (updated %s)\n",
						name, sm.EntryCount, sm.LastUpdated.Format("2006-01-02 15:04"))
				}
			}

			return nil
		},
	}
}

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push <oci-reference>",
		Short: "Push the cached database to an OCI registry",
		Long: `Push the local vulnerability database to an OCI registry as an artifact.

Examples:
  taco-vulndb push ghcr.io/myorg/taco-vulndb:latest
  taco-vulndb push ghcr.io/myorg/taco-vulndb:2024-01-15`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := args[0]

			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			if !quiet {
				slog.Info("pushing vulnerability database...", "ref", ref)
			}

			if err := vulndb.PushOCI(cache, ref); err != nil {
				return fmt.Errorf("pushing database: %w", err)
			}

			if !quiet {
				slog.Info("database pushed successfully", "ref", ref)
			}

			return nil
		},
	}
}

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull <oci-reference>",
		Short: "Pull a database from an OCI registry into the local cache",
		Long: `Pull a vulnerability database OCI artifact from a registry.

Examples:
  taco-vulndb pull ghcr.io/myorg/taco-vulndb:latest`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := args[0]

			cache, err := vulndb.NewCache()
			if err != nil {
				return fmt.Errorf("initializing cache: %w", err)
			}

			if !quiet {
				slog.Info("pulling vulnerability database...", "ref", ref)
			}

			if err := vulndb.PullOCI(cache, ref); err != nil {
				return fmt.Errorf("pulling database: %w", err)
			}

			meta, _ := cache.ReadMeta()
			if !quiet && meta != nil {
				slog.Info("database pulled successfully",
					"ref", ref,
					"entries", meta.EntryCount,
				)
			}

			return nil
		},
	}
}
