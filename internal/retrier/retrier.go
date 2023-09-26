package retrier

import (
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tanya-mtv/metricsservice/internal/constants"
)

type Retrier struct {
	Retries []time.Duration
}

func NewRetrier() *Retrier {
	return &Retrier{
		Retries: []time.Duration{0, constants.RetryWaitMin, constants.RetryMedium, constants.RetryWaitMax},
	}
}

func (r *Retrier) Next(err error, timeint time.Duration) bool {

	if r.HaveToRetry(err) {
		time.Sleep(timeint)
		return true
	}
	return false
}

func (r *Retrier) HaveToRetry(err error) bool {
	var pge *pgconn.PgError
	var ose *os.PathError

	if errors.As(err, &pge) {
		return true
	}
	if errors.As(err, &ose) {
		return true
	}
	return false
}
