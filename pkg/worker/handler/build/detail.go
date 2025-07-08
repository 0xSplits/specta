package build

type detail struct {
	// repo is the repository name that we are fetching build information for from
	// the Github API.
	repo string
	// label is the metric label used to instrument this project in Grafana.
	label string
	// check is the name of the workflow within this repository running all
	// relevant verification steps for this project.
	check string
	// image is the name of the workflow within this repository pushing the docker
	// image during the release process.
	image string
}

// detail returns a hard coded list of build information applicable for every
// environment.
func (h *Handler) detail() ([]detail, error) {
	var det []detail
	{
		det = []detail{
			{
				repo:  "specta",
				label: "specta",
				check: "go-build",
				image: "docker-release",
			},
			{
				repo:  "splits",
				label: "server",
				check: "typescript-server",
				image: "docker-push",
			},
		}
	}

	return det, nil
}
