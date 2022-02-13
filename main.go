package main

import (
	"bytes"
	"os"
	"strconv"
	"text/template"

	"github.com/deifyed/deviation-service/pkg/events"
	"github.com/deifyed/deviation-service/pkg/notification"
	"github.com/deifyed/deviation-service/pkg/urlgenerator"
)

func main() {
	const (
		timeType  = "time"
		valueType = "value"
	)

	var (
		deviation                = deviationHandler{}
		eventsClient             events.Client
		timeDeviationCalculator  movingAverageCalculator
		valueDeviationCalculator movingAverageCalculator
		timeThreshold, _         = strconv.ParseFloat(os.Getenv("TIME_THRESHOLD"), 64)
		valueThreshold, _        = strconv.ParseFloat(os.Getenv("VALUE_THRESHOLD"), 64)
	)

	for {
		event := eventsClient.Receive()

		timeDeviationCalculator.Update(event.Timestamp.UnixMilli())
		valueDeviationCalculator.Update(event.Value)

		deviation.Handle(
			timeThreshold,
			timeDeviationCalculator.StandardDeviation(),
			timeType,
			event.ID,
		)

		deviation.Handle(
			valueThreshold,
			valueDeviationCalculator.StandardDeviation(),
			valueType,
			event.ID,
		)
	}
}

const alertTemplate = `
Found {{ .Type }} deviation above threshold!

See [Grafana]({{ .GrafanaURL }}) for more details.
Manage data in [Chronograf]({{ .ChronografURL }}).

Godspeed.
`

type alertTemplateOpts struct {
	Type          string
	GrafanaURL    string
	ChronografURL string
}

type deviationHandler struct {
	notifier     notification.Client
	urlGenerator urlgenerator.Client
}

func (d deviationHandler) Handle(threshold, deviation float64, deviationType, eventID string) {
	if deviation < threshold {
		return
	}

	buf := bytes.Buffer{}

	t, _ := template.New("").Parse(alertTemplate)

	_ = t.Execute(&buf, alertTemplateOpts{
		Type:          deviationType,
		GrafanaURL:    d.urlGenerator.GenerateGrafanaURLForID(eventID),
		ChronografURL: d.urlGenerator.GenerateChronografURLForID(eventID),
	})

	d.notifier.Send("Warning!", buf.String())
}
