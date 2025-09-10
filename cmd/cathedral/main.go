package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oboro/cathedral/pkg/config"
	"github.com/oboro/cathedral/pkg/store"
	"github.com/oboro/cathedral/pkg/session"
	"github.com/oboro/cathedral/pkg/memory"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	configPath string
	verbose    bool
	
	// Command-specific flags
	storePath      string
	storeName      string
	sessionID      string
	templatePath   string
	indexPath      string
	getPromptOnly  bool
	noTags         bool
	compression    float64
	compressionProfile string
	dateInput      string
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

	// Session management commands
	initSessionCmd := &cobra.Command{
		Use:   "init-episodic-session",
		Short: "Initialize a new episodic session",
		Args:  cobra.NoArgs,
		RunE:  runInitEpisodicSession,
	}
	initSessionCmd.Flags().StringVar(&dateInput, "date", "", "Date/time for the session (YYYY-MM-DD, YYYYMMDD, or unix timestamp)")

	// Import command
	importCmd := &cobra.Command{
		Use:   "import-messages FILES...",
		Short: "Import messages from Hinata format",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runImportMessages,
	}
	importCmd.Flags().StringVar(&sessionID, "session", "", "Existing session to append to (format: YYYYMMDD/SESSION_ID)")

	// Memory writing command
	writeMemoryCmd := &cobra.Command{
		Use:   "write-memory",
		Short: "Generate memory writing prompt for a session",
		Args:  cobra.NoArgs,
		RunE:  runWriteMemory,
	}
	writeMemoryCmd.Flags().StringVar(&sessionID, "session", "", "Session to process (default: latest)")
	writeMemoryCmd.Flags().StringVar(&templatePath, "template", "", "Template file to use")
	writeMemoryCmd.Flags().StringVar(&indexPath, "index", "", "Index file to use")
	writeMemoryCmd.Flags().BoolVar(&getPromptOnly, "get-prompt", false, "Only output the prompt without submitting to LLM")
	writeMemoryCmd.Flags().Float64Var(&compression, "compression", 0.5, "Compression ratio (0.0-1.0)")
	writeMemoryCmd.Flags().StringVar(&compressionProfile, "compression-profile", "", "Use predefined compression profile")

	// Conversation init command
	initConvCmd := &cobra.Command{
		Use:   "init-session",
		Short: "Initialize a new conversation session with memory index",
		Args:  cobra.NoArgs,
		RunE:  runInitConversation,
	}
	initConvCmd.Flags().StringVar(&templatePath, "template", "", "Template file to use")
	initConvCmd.Flags().BoolVar(&getPromptOnly, "get-prompt", false, "Only output the prompt without creating a session")

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
		initSessionCmd,
		importCmd,
		writeMemoryCmd,
		initConvCmd,
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

func runInitEpisodicSession(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	
	sessMgr := session.NewManager(cfg)
	sessionPath, err := sessMgr.InitEpisodicSession(dateInput)
	if err != nil {
		return err
	}
	
	fmt.Println(sessionPath)
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

func runWriteMemory(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	
	// Handle compression profile
	if compressionProfile != "" {
		if ratio, ok := memory.CompressionProfiles[compressionProfile]; ok {
			compression = ratio
		} else {
			return fmt.Errorf("unknown compression profile: %s", compressionProfile)
		}
	}
	
	writer := memory.NewWriter(cfg)
	return writer.WriteMemory(sessionID, templatePath, indexPath, getPromptOnly, compression)
}

func runInitConversation(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	
	conv := session.NewConversationInitializer(cfg)
	return conv.InitSession(templatePath, getPromptOnly)
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