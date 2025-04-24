package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    vmware "SDVCLI/Auth"
)

var StopCmd = &cobra.Command{
    Use:   "stop",
    Short: "Arrêter une VM",
    Run: func(cmd *cobra.Command, args []string) {
        config, err := vmware.CheckConfiguration()
        if err != nil {
            fmt.Println("Aucune configuration trouvée")
            return
        }
        if len(args) < 1 {
            fmt.Println("Veuillez spécifier l'ID de la VM à arrêter")
            return
        }
        vmID := args[0]
        err = vmware.StopVM(config, vmID)
        if err != nil {
            fmt.Println("Impossible d'arrêter la VM")
            fmt.Println(err)
            return
        }
        fmt.Printf("La VM avec l'ID %s a été arrêtée avec succès.\n", vmID)
    },
}

func init() {
    vmCmd.AddCommand(StopCmd)
}
