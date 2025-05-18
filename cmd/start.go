package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    vmware "SDVCLI/Auth"
)

var StartCmd = &cobra.Command{
    Use:   "start",
    Short: "Démarrer une VM",
    Run: func(cmd *cobra.Command, args []string) {
        config, err := vmware.CheckConfiguration()
        if err != nil {
            fmt.Println("Aucune configuration trouvée")
            return
        }
        if len(args) < 1 {
            fmt.Println("Veuillez spécifier l'ID de la VM à démarrer")
            return
        }
        vmID := args[0]
        err = vmware.StartVM(config, vmID)
        if err != nil {
            fmt.Println("Impossible de démarrer la VM")
            fmt.Println(err)
            return
        }
        fmt.Printf("La VM avec %s a été démarrée avec succès.\n", vmID)
    },
}

func init() {
   ServeurCmd.AddCommand(StartCmd)
}
