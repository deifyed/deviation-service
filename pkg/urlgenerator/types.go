package urlgenerator

type Client interface {
	GenerateGrafanaURLForID(id string) string
	GenerateChronografURLForID(id string) string
}
