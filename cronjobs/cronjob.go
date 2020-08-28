package cronjobs

import (
	"context"
	"github.com/zdarovich/win-lose-api/database"
	"github.com/zdarovich/win-lose-api/models"
	"log"
	"strconv"
	"sync"
	"time"
)

const IntervalPeriod = 10 * time.Second

type jobTicker struct {
	ticker *time.Ticker
}

var (
	cron *jobTicker
	cronjobSync sync.Once
)


func GetInstance() *jobTicker {
	cronjobSync.Do(func() {
		cron = &jobTicker{ticker:time.NewTicker(IntervalPeriod)}
	})
	return cron
}

// cronjob scheduler
func (jobTicker *jobTicker) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-jobTicker.ticker.C:
				jobTicker.task(ctx)
			case <- ctx.Done():
				jobTicker.ticker.Stop()
				return
			}
		}
	}()
}

// cronjob task that rollbacks last 10 odd transactions
func (jobTicker *jobTicker) task(ctx context.Context) {
	log.Println("cron job started")
	defer log.Println("cron job finished")

	tx := database.GetDB().Begin()
	if tx.Error != nil {
		log.Println(tx.Error)
		return
	}
	rows, err := tx.Raw("select * from transactions where ( id % 2 ) != 0 and canceled = false limit 10").Rows()
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var trxs []*models.Transaction
	for rows.Next() {
		var trx models.Transaction
		err := tx.ScanRows(rows, &trx)
		if err != nil {
			log.Println(err)
			return
		}
		trxs = append(trxs, &trx)
	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	if len(trxs) == 0 {
		return
	}
	var rollBackAmount float64
	for _, trx := range trxs {
		if err := tx.Model(trx).Updates(map[string]interface{}{"canceled": true}).Error; err != nil {
			log.Println(err)
			tx.Rollback()
			return
		}
		f, err := strconv.ParseFloat(trx.Amount, 64)
		if err != nil {
			log.Println(err)
			return
		}

		// it sums up all transaction amounts and reverts their sign
		if trx.State == "lost" {
			rollBackAmount += f
		} else {
			rollBackAmount += -f
		}
		log.Println("transaction is updated by cronjob")
		log.Printf("%+v \n", trx)
	}

	if tx.Error != nil {
		log.Println(tx.Error)
		return
	}
	var user models.User
	if err := tx.First(&user).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}
	if err := tx.Model(&user).Update("balance", user.Balance + rollBackAmount).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}
	log.Println("user balance is updated by cronjob")
	log.Printf("%+v \n", user)
	tx.Commit()
}


