package pgx

import (
	"context"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"payment_processing_system/internal/domain"
	"payment_processing_system/pkg/logger"
	"testing"
)

type BalanceStorageTestSuite struct {
	suite.Suite
	pool pgxmock.PgxPoolIface
	logg *logger.Logger
}

func (suite *BalanceStorageTestSuite) SetupTest() {
	var err error
	var l *zap.Logger
	l, err = zap.NewProduction()
	suite.logg = logger.NewLogger(l)
	if err != nil {
		suite.Fail(err.Error())
	}

	suite.pool, err = pgxmock.NewPool()
	if err != nil {
		suite.Fail(err.Error())
	}
}

func (suite *BalanceStorageTestSuite) TearDownTest() {
	suite.pool.Close()
}

func (suite *BalanceStorageTestSuite) TestIncreaseAmount_ExecErr() {
	var expectedErr = errors.New("exec error")

	ctx := context.Background()
	var id int64 = 1
	amount := decimal.NewFromInt(3)
	bs := NewBalanceStorage(suite.pool, suite.logg)

	suite.pool.ExpectExec(`^INSERT INTO public\.balance \(id,amount\) VALUES (.+,.+) ON CONFLICT DO UPDATE SET amount = amount \+ .+$`).
		WithArgs(id, amount, amount).WillReturnError(expectedErr)

	err := bs.IncreaseAmount(ctx, id, amount)

	suite.Equal(expectedErr, err)
}

func (suite *BalanceStorageTestSuite) TestIncreaseAmount_BalanceWasNotIncreasedErr() {
	var expectedErr = domain.BalanceWasNotIncreasedErr

	ctx := context.Background()
	var id int64 = 1
	amount := decimal.NewFromInt(3)
	bs := NewBalanceStorage(suite.pool, suite.logg)

	suite.pool.ExpectExec(`^INSERT INTO public\.balance \(id,amount\) VALUES (.+,.+) ON CONFLICT DO UPDATE SET amount = amount \+ .+$`).
		WithArgs(id, amount, amount).WillReturnResult(pgxmock.NewResult("INSERT", 0))

	err := bs.IncreaseAmount(ctx, id, amount)

	suite.Equal(expectedErr, err)
}

func (suite *BalanceStorageTestSuite) TestIncreaseAmount_Success() {
	var expectedErr error

	ctx := context.Background()
	var id int64 = 1
	amount := decimal.NewFromInt(3)
	bs := NewBalanceStorage(suite.pool, suite.logg)

	suite.pool.ExpectExec(`^INSERT INTO public\.balance \(id,amount\) VALUES (.+,.+) ON CONFLICT DO UPDATE SET amount = amount \+ .+$`).
		WithArgs(id, amount, amount).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := bs.IncreaseAmount(ctx, id, amount)

	suite.Equal(expectedErr, err)
}

func TestBalanceStorageTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceStorageTestSuite))
}
