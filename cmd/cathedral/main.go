package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/veilm/cathedral/pkg/config"
	"github.com/veilm/cathedral/pkg/memory"
	"github.com/veilm/cathedral/pkg/session"
	"github.com/veilm/cathedral/pkg/store"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	configPath string
	verbose    bool

	// Command-specific flags
	storePath          string
	storeName          string
	sessionID          string
	templatePath       string
	indexPath          string
	getPromptOnly      bool
	noTags             bool
	compression        float64
	compressionProfile string
	dateInput          string
	sleepSessionDir    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cathedral",
		Short: "Cathedral memory store management",
		Long: `Cathedral is a memory system for managing conversation histories
and knowledge bases with episodic and semantic memory structures.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Config file path (default: $XDG_CONFIG_HOME/cathedral/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Store management commands
	createStoreCmd := &cobra.Command{
		Use:   "create-store NAME [PATH]",
		Short: "Create a new memory store",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runCreateStore,
	}

	linkStoreCmd := &cobra.Command{
		Use:   "link-store PATH",
		Short: "Link an existing directory as a memory store",
		Args:  cobra.ExactArgs(1),
		RunE:  runLinkStore,
	}
	linkStoreCmd.Flags().StringVar(&storeName, "name", "", "Name for the memory store (default: directory basename)")

	listStoresCmd := &cobra.Command{
		Use:   "list-stores",
		Short: "List all memory stores",
		Args:  cobra.NoArgs,
		RunE:  runListStores,
	}

	switchStoreCmd := &cobra.Command{
		Use:   "switch-store NAME",
		Short: "Switch to a different memory store",
		Args:  cobra.ExactArgs(1),
		RunE:  runSwitchStore,
	}

	unlinkStoreCmd := &cobra.Command{
		Use:   "unlink-store NAME_OR_PATH",
		Short: "Unlink a memory store from config (does not delete files)",
		Args:  cobra.ExactArgs(1),
		RunE:  runUnlinkStore,
	}

	showActiveCmd := &cobra.Command{
		Use:   "show-active-store",
		Short: "Show the currently active store",
		Args:  cobra.NoArgs,
		RunE:  runShowActive,
	}

	// Memory episode management commands
	initMemoryEpisodeCmd := &cobra.Command{
		Use:   "init-memory-episode",
		Short: "Initialize a new memory episode for storing conversations",
		Args:  cobra.NoArgs,
		RunE:  runInitMemoryEpisode,
	}
	initMemoryEpisodeCmd.Flags().StringVar(&dateInput, "date", "", "Date/time for the episode (YYYY-MM-DD, YYYYMMDD, or unix timestamp)")

	// Import command
	importCmd := &cobra.Command{
		Use:   "import-messages FILES...",
		Short: "Import messages from Hinata format",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runImportMessages,
	}
	importCmd.Flags().StringVar(&sessionID, "session", "", "Existing session to append to (format: YYYYMMDD/SESSION_ID)")

	// Consolidation planning command
	planConsolidationCmd := &cobra.Command{
		Use:   "plan-consolidation",
		Short: "Generate consolidation planning prompt for a session",
		Args:  cobra.NoArgs,
		RunE:  runPlanConsolidation,
	}
	planConsolidationCmd.Flags().StringVar(&sessionID, "session", "", "Session to process (default: latest)")
	planConsolidationCmd.Flags().StringVar(&templatePath, "template", "", "Template file to use")
	planConsolidationCmd.Flags().StringVar(&indexPath, "index", "", "Index file to use")
	planConsolidationCmd.Flags().BoolVar(&getPromptOnly, "prepare-only", false, "Only prepare the conversation without running agent loop")
	planConsolidationCmd.Flags().Float64Var(&compression, "compression", 0.5, "Compression ratio (0.0-1.0)")
	planConsolidationCmd.Flags().StringVar(&compressionProfile, "compression-profile", "", "Use predefined compression profile")

	// Consolidation execution command
	executeConsolidationCmd := &cobra.Command{
		Use:   "execute-consolidation",
		Short: "Execute a consolidation plan from a sleep session",
		Long: `Execute all operations from a consolidation plan.
Use --session to specify which sleep session to execute, or omit to use the latest.`,
		Args: cobra.NoArgs,
		RunE: runExecuteConsolidation,
	}
	executeConsolidationCmd.Flags().StringVar(&sleepSessionDir, "session", "", "Sleep session (timestamp or full path, default: latest)")

	// Conversation start command
	startConvCmd := &cobra.Command{
		Use:   "start-conversation",
		Short: "Start a new conversation with memory context injected",
		Args:  cobra.NoArgs,
		RunE:  runStartConversation,
	}
	startConvCmd.Flags().StringVar(&templatePath, "template", "", "Template file to use")
	startConvCmd.Flags().BoolVar(&getPromptOnly, "get-prompt", false, "Only output the prompt without creating a hnt-chat session")

	// Health check command
	healthCmd := &cobra.Command{
		Use:   "check-health [FILES...]",
		Short: "Check health of memory node files by validating [[links]]",
		Args:  cobra.ArbitraryArgs,
		RunE:  runHealthCheck,
	}

	// Read command
	readCmd := &cobra.Command{
		Use:     "read FILES...",
		Aliases: []string{"read-node"},
		Short:   "Read memory node files to stdout",
		Args:    cobra.MinimumNArgs(1),
		RunE:    runReadNode,
	}
	readCmd.Flags().BoolVar(&noTags, "no-tags", false, "Don't wrap output in XML tags")

	// Add all commands to root
	rootCmd.AddCommand(
		createStoreCmd,
		linkStoreCmd,
		listStoresCmd,
		switchStoreCmd,
		unlinkStoreCmd,
		showActiveCmd,
		initMemoryEpisodeCmd,
		importCmd,
		planConsolidationCmd,
		executeConsolidationCmd,
		startConvCmd,
		healthCmd,
		readCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Command implementations (stubs for now)
func runCreateStore(cmd *cobra.Command, args []string) error {
	name := args[0]
	path := ""
	if len(args) > 1 {
		path = args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	mgr := store.NewManager(cfg)
	return mgr.CreateStore(name, path)
}

func runLinkStore(cmd *cobra.Command, args []string) error {
	path := args[0]

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if storeName == "" {
		storeName = filepath.Base(path)
	}

	mgr := store.NewManager(cfg)
	return mgr.LinkStore(storeName, path)
}

func runListStores(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	mgr := store.NewManager(cfg)
	return mgr.ListStores()
}

func runSwitchStore(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	mgr := store.NewManager(cfg)
	return mgr.SwitchStore(name)
}

func runUnlinkStore(cmd *cobra.Command, args []string) error {
	nameOrPath := args[0]

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	mgr := store.NewManager(cfg)
	return mgr.UnlinkStore(nameOrPath)
}

func runShowActive(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	mgr := store.NewManager(cfg)
	return mgr.ShowActive()
}

func runInitMemoryEpisode(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	sessMgr := session.NewManager(cfg)
	episodePath, err := sessMgr.InitMemoryEpisode(dateInput)
	if err != nil {
		return err
	}

	fmt.Println(episodePath)
	return nil
}

func runImportMessages(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	importer := session.NewImporter(cfg)
	return importer.ImportMessages(args, sessionID)
}

func runStartConversation(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	conv := session.NewConversationStarter(cfg)
	return conv.StartConversation(templatePath, getPromptOnly)
}

func runHealthCheck(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	checker := memory.NewHealthChecker(cfg)
	return checker.CheckHealth(args)
}

func runReadNode(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	reader := memory.NewNodeReader(cfg)
	return reader.ReadNodes(args, noTags)
}

func runPlanConsolidation(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	// Handle compression profile
	if compressionProfile != "" {
		if ratio, ok := config.CompressionProfiles[compressionProfile]; ok {
			compression = ratio
		} else {
			return fmt.Errorf("unknown compression profile: %s", compressionProfile)
		}
	}

	planner := memory.NewPlanner(cfg)
	return planner.PlanConsolidation(sessionID, templatePath, indexPath, getPromptOnly, compression)
}

func runExecuteConsolidation(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	activeStore := cfg.GetActiveStorePath()
	if activeStore == "" {
		return fmt.Errorf("no active memory store. Create one with 'cathedral create-store <name>'")
	}

	sleepDir := filepath.Join(activeStore, "sleep")

	// Resolve sleep session directory
	var sleepSessionPath string
	if sleepSessionDir != "" {
		// Check if it's a full path
		if filepath.IsAbs(sleepSessionDir) {
			sleepSessionPath = sleepSessionDir
		} else {
			// Assume it's a timestamp
			sleepSessionPath = filepath.Join(sleepDir, sleepSessionDir)
		}
	} else {
		// Find latest sleep session
		entries, err := os.ReadDir(sleepDir)
		if err != nil {
			return fmt.Errorf("failed to read sleep directory: %w", err)
		}

		// Find the most recent directory (highest timestamp)
		var latestTimestamp string
		for _, entry := range entries {
			if entry.IsDir() && entry.Name() > latestTimestamp {
				latestTimestamp = entry.Name()
			}
		}

		if latestTimestamp == "" {
			return fmt.Errorf("no sleep sessions found in %s", sleepDir)
		}

		sleepSessionPath = filepath.Join(sleepDir, latestTimestamp)
		fmt.Printf("Using latest sleep session: %s\n", latestTimestamp)
	}

	// Verify the directory exists
	if _, err := os.Stat(sleepSessionPath); os.IsNotExist(err) {
		return fmt.Errorf("sleep session directory not found: %s", sleepSessionPath)
	}

	executor := memory.NewExecutor(cfg)
	return executor.ExecuteConsolidation(sleepSessionPath)
}
