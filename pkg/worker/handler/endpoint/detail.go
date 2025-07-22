package endpoint

type detail struct {
	// lab is the metric label used to instrument this service endpoint in
	// Grafana.
	lab string
	// url is the list of HTTPS endpoints to verify for this service.
	url []string
}

// detail returns a hard coded list of environment specific service endpoints to
// verify.
func (h *Handler) detail() ([]detail, error) {
	var det map[string][]detail
	{
		det = map[string][]detail{
			"testing": {
				{
					lab: "explorer",
					url: []string{"https://explorer.testing.splits.org"},
				},
				{
					lab: "server",
					url: []string{"https://server.testing.splits.org/metrics"},
				},
				{
					lab: "specta",
					url: []string{"https://specta.testing.splits.org/metrics"},
				},
				{
					lab: "teams",
					url: []string{"https://teams.testing.splits.org"},
				},
			},
			"staging": {
				{
					lab: "explorer",
					url: []string{"https://explorer.staging.splits.org"},
				},
				{
					lab: "server",
					url: []string{"https://server.staging.splits.org/metrics"},
				},
				{
					lab: "specta",
					url: []string{"https://specta.staging.splits.org/metrics"},
				},
				{
					lab: "teams",
					url: []string{"https://teams.staging.splits.org"},
				},
			},
			"production": {
				{
					lab: "explorer",
					url: []string{"https://explorer.production.splits.org", "https://app.splits.org"},
				},
				{
					lab: "server",
					url: []string{"https://server.production.splits.org/metrics", "https://api.splits.org/metrics"},
				},
				{
					lab: "specta",
					url: []string{"https://specta.production.splits.org/metrics"},
				},
				{
					lab: "teams",
					url: []string{"https://teams.production.splits.org", "https://teams.splits.org"},
				},
			},
		}
	}

	return det[h.env.Environment], nil
}
