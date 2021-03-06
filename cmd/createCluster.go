package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Krucible/krucible-go-client/krucible"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type durationFlag struct {
	duration *int
}

var _ pflag.Value = &durationFlag{}

func (df *durationFlag) String() string {
	return ""
}

func (df *durationFlag) Set(val string) error {
	if val == "permanent" {
		df.duration = nil
		return nil
	}

	i, err := strconv.Atoi(val)
	if err != nil || i < 1 || i > 6 {
		return fmt.Errorf(`Cluster Duration must be an integer between 1 and 6 or "permanent"`)
	}

	df.duration = &i
	return nil
}

func (df *durationFlag) Type() string {
	return `integer or "permanent"`
}

var DisplayName string
var ClusterDuration durationFlag
var ConfigureKubectlFlag bool

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a new cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClientOrDie()
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Creating cluster... "
		s.Start()
		cluster, _, err := client.CreateCluster(krucible.CreateClusterConfig{
			DisplayName:     DisplayName,
			DurationInHours: ClusterDuration.duration,
		})
		s.Stop()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, cluster.ID)
		fmt.Fprintln(os.Stderr, "")

		if ConfigureKubectlFlag {
			configureKubectl(
				cluster.ID,
				cluster.ConnectionDetails.Server,
				cluster.ConnectionDetails.CertificateAuthority,
				cluster.ConnectionDetails.ClusterAuthToken,
			)
		}
	},
}

func init() {
	createCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().BoolVarP(&ConfigureKubectlFlag, "configure-kubectl", "k", false, "Configure kubectl context for connection to your cluster")
	clusterCmd.Flags().StringVarP(&DisplayName, "name", "n", "", "Desired display name for the cluster")
	clusterCmd.Flags().VarP(&ClusterDuration, "duration", "d", "The amount of time the cluster should persist for")
	clusterCmd.MarkFlagRequired("name")
	clusterCmd.MarkFlagRequired("duration")
}
