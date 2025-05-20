package cmd

import (
	"fmt"
	"os"
	vmware "SDVCLI/Auth"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [vmID]",
	Short: "Supprimer une vm",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmID := args[0]

		config, err := vmware.CheckConfiguration()
		if err != nil {
			fmt.Println("Aucune configuration trouvée")
			return
		}

		err = vmware.DeleteVM(config, vmID)
		if err != nil {
			fmt.Printf("Erreur lors de la suppression de la VM : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("VM supprimée avec succès.")
	},
}

func init() {
	ServeurCmd.AddCommand(DeleteCmd)
}
