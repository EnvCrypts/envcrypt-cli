package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for envcrypt.

To load completions:

Bash:
  $ source <(envcrypt completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ envcrypt completion bash > /etc/bash_completion.d/envcrypt
  # macOS:
  $ envcrypt completion bash > $(brew --prefix)/etc/bash_completion.d/envcrypt

Zsh:
  $ source <(envcrypt completion zsh)
  # To load completions for each session, execute once:
  $ envcrypt completion zsh > "${fpath[1]}/_envcrypt"

Fish:
  $ envcrypt completion fish | source
  # To load completions for each session, execute once:
  $ envcrypt completion fish > ~/.config/fish/completions/envcrypt.fish

PowerShell:
  PS> envcrypt completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, execute once:
  PS> envcrypt completion powershell > envcrypt.ps1`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
