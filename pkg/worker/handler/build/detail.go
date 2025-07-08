package build

type detail struct {
	// name is the repository name that we are instrumenting.
	name string
	// label is the metric label used to instrument this project.
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
				name:  "specta",
				label: "specta",
				check: "go-build",
				image: "docker-release",
			},
			{
				name:  "splits",
				label: "server",
				check: "typescript-server",
				image: "docker-push",
			},
		}
	}

	return det, nil
}
