package cmd

import (
	"fmt"

	vmware "SDVCLI/Auth"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "list",
	Short: "Liste les vm",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := vmware.CheckConfiguration()
		if err != nil {
			fmt.Println("Aucune configuration trouvée")
			return
		}
		VMList, err := vmware.GetVMlist(config)
		if err != nil {
			fmt.Println("Impossible de recupérer la liste des vm")
			fmt.Println(err)
			return
		}
		vmware.AffichageVM(VMList)
	},
}

func init() {
	ServeurCmd.AddCommand(StartCmd)
}
