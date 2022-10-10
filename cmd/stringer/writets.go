package stringer

import (
	"context"
	"fmt"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"

	"github.com/spf13/cobra"
)

var writeTsCmd = &cobra.Command{
	Use:     "writets",
	Aliases: []string{"w"},
	Short:   "Writes a time series",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Creates a client.
		client, err := monitoring.NewMetricClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		// Sets your Google Cloud Platform project ID.
		// Run `gcloud auth application-default login` first to configure client authentication/authorization
		projectID := "slo-generator-demo"

		// Prepares an individual data point
		dataPoint := &monitoringpb.Point{
			Interval: &monitoringpb.TimeInterval{
				EndTime: &googlepb.Timestamp{
					Seconds: time.Now().Unix(),
				},
			},
			Value: &monitoringpb.TypedValue{
				Value: &monitoringpb.TypedValue_DoubleValue{
					DoubleValue: 123.45,
				},
			},
		}

		// Writes time series data.
		if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
			Name: fmt.Sprintf("projects/%s", projectID),
			TimeSeries: []*monitoringpb.TimeSeries{
				{
					Metric: &metricpb.Metric{
						Type: "custom.googleapis.com/stores/daily_sales",
						Labels: map[string]string{
							"store_id": "Pittsburg",
						},
					},
					Resource: &monitoredrespb.MonitoredResource{
						Type: "global",
						Labels: map[string]string{
							"project_id": projectID,
						},
					},
					Points: []*monitoringpb.Point{
						dataPoint,
					},
				},
			},
		}); err != nil {
			log.Fatalf("Failed to write time series data: %v", err)
		}

		// Closes the client and flushes the data to Stackdriver.
		if err := client.Close(); err != nil {
			log.Fatalf("Failed to close client: %v", err)
		}

		fmt.Printf("Done writing time series data.\n")
	},
}

func init() {
	rootCmd.AddCommand(writeTsCmd)
}
