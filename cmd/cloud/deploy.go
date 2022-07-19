package cloud

import (
	"fmt"

	cloud "github.com/astronomer/astro-cli/cloud/deploy"
	"github.com/astronomer/astro-cli/cmd/utils"
	"github.com/astronomer/astro-cli/config"
	"github.com/astronomer/astro-cli/pkg/git"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	forceDeploy      bool
	forcePrompt      bool
	saveDeployConfig bool
	pytest           bool
	deployExample    = `
Specify the ID of the Deployment on Astronomer you would like to deploy this project to:

  $ astro deploy <deployment ID>

Menu will be presented if you do not specify a deployment ID:

  $ astro deploy
`

	deployImage      = cloud.Deploy
	ensureProjectDir = utils.EnsureProjectDir
)

var (
	pytestFile string
	envFile    string
	imageName  string
)

const (
	registryUncommitedChangesMsg = "Project directory has uncommitted changes, use `astro deploy [deployment-id] -f` to force deploy."
)

func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy DEPLOYMENT-ID",
		Short:   "Deploy your project to a Deployment on Astro",
		Long:    "Deploy your project to a Deployment on Astro. This command bundles your project files into a Docker image and pushes that Docker image to Astronomer. It does not include any metadata associated with your local Airflow environment.",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: ensureProjectDir,
		RunE:    deploy,
		Example: deployExample,
	}
	cmd.Flags().BoolVarP(&forceDeploy, "force", "f", false, "Force deploy even if project contains errors or uncommitted changes")
	cmd.Flags().BoolVarP(&forcePrompt, "prompt", "p", false, "Force prompt to choose target deployment")
	cmd.Flags().BoolVarP(&saveDeployConfig, "save", "s", false, "Save deployment in config for future deploys")
	cmd.Flags().StringVar(&workspaceID, "workspace-id", "", "Workspace for your Deployment")
	cmd.Flags().BoolVar(&pytest, "pytest", false, "Deploy code to Astro only if the specified Pytests are passed")
	cmd.Flags().StringVarP(&envFile, "env", "e", ".env", "Location of file containing environment variables for Pytests")
	cmd.Flags().StringVarP(&pytestFile, "test", "t", "", "Location of Pytests or specific Pytest file. All Pytest files must be located in the tests directory")
	cmd.Flags().StringVarP(&imageName, "image-name", "i", "", "Name of a custom image to deploy")
	return cmd
}

func deploy(cmd *cobra.Command, args []string) error {
	deploymentID := ""
	ws := ""

	// Get release name from args, if passed
	if len(args) > 0 {
		deploymentID = args[0]
	}

	if deploymentID == "" || forcePrompt {
		var err error
		ws, err = coalesceWorkspace()
		if err != nil {
			return errors.Wrap(err, "failed to find a valid workspace")
		}
	}

	// Save release name in config if specified
	if len(deploymentID) > 0 && saveDeployConfig {
		err := config.CFG.ProjectDeployment.SetProjectString(deploymentID)
		if err != nil {
			return nil
		}
	}

	if git.HasUncommittedChanges() && !forceDeploy {
		fmt.Println(registryUncommitedChangesMsg)
		return nil
	}

	if pytest && pytestFile == "" {
		pytestFile = "all-tests"
	}

	if !pytest && !forceDeploy {
		pytestFile = "parse"
	}

	// Silence Usage as we have now validated command input
	cmd.SilenceUsage = true

	return deployImage(config.WorkingPath, deploymentID, ws, pytestFile, envFile, imageName, forcePrompt, astroClient)
}