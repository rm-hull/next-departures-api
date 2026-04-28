package internal

import (
	"log"

	"github.com/rm-hull/next-departures-api/internal/models"
	"github.com/robfig/cron/v3"
)

const CRON_SCHEDULE_NAPTAN = "0 */19 * * *" // Every 19 hours

func StartCron(repo NaptanRepository) (*cron.Cron, error) {

	log.Printf("Starting CRON job to update NaPTAN datasets (schedule: %s)", CRON_SCHEDULE_NAPTAN)

	c := cron.New()
	if _, err := c.AddFunc(CRON_SCHEDULE_NAPTAN, func() {
		err := TransientDownload(models.NAPTAN_CSV_URL, repo.ImportCSV)
		if err != nil {
			log.Printf("Error importing download NaPTAN dataset: %v", err)
		}
	}); err != nil {
		return nil, err
	}

	c.Start()
	return c, nil
}
