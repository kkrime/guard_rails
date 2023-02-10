package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"

	"guard_rails/db"
	"guard_rails/db/mock"
)

type ScanServiceTestSuite struct {
	suite.Suite
}

func (s *ScanServiceTestSuite) TestupdateTransactionLimts() {

	tests := []struct {
		name                                  string
		scanDBmoc                             db.ScanDb
		mockTimeInterface                     transactionLimitsDateTimeHelper.TimeHelper
		mockLimitsRepository                  repository.LimitsRepositoryProvider
		mockMonoLithHttpInterface             client.MonolitClientProvider
		mockTransactionLimitsUpdateRepository repository.TransactionLimitsUpdateRepositoryProvider
		out                                   []*service.UpdateTransactionLimitsResult
		expectedLog                           []string
		err                                   error
	}{
		{
			name: "Adding New Limit with -1 value",

			transactionLimits: []*model.TransactionLimit{
				{
					LimitType:    model.Daily,
					Amount:       -1,
					CurrencyCode: "SGD",
				},
			},

			mockLimitsRepository: func() repository.LimitsRepositoryProvider {
				lrMock := limitsRepositoryMock.NewLimitsRepositoryProvider(s.T())

				lrMock.On(
					"CreateEntity",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
						) error {
							return nil
						},
					).
					Once()

				return lrMock
			}(),
		},
		{
			name: "Adding New Limit with -2 value (delete)",

			transactionLimits: []*model.TransactionLimit{
				{
					LimitType:    model.Daily,
					Amount:       -2,
					CurrencyCode: "SGD",
				},
			},

			mockLimitsRepository: func() repository.LimitsRepositoryProvider {
				lrMock := limitsRepositoryMock.NewLimitsRepositoryProvider(s.T())

				lrMock.On(
					"CreateEntity",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
						) error {
							return nil
						},
					).
					Once()

				lrMock.On(
					"DeleteTransactionLimitType",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entityId string,
							limitType model.LimitType,
						) error {
							return nil
						},
					).
					Once()

				return lrMock
			}(),

			expectedLog: []string{
				fmt.Sprintf(literal.DeleteingTransactionLimit, entity.EntityType, entity.EntityId, model.Daily),
			},
		},
		{
			name: "Adding New Limit",

			out: []*service.UpdateTransactionLimitsResult{{
				TransactionLimit: &model.TransactionLimit{
					LimitType:    model.Daily,
					CurrencyCode: "SGD",
					Amount:       33.33,
				},
			}},

			transactionLimits: []*model.TransactionLimit{
				{
					LimitType:    model.Daily,
					Amount:       33.33,
					CurrencyCode: "SGD",
				},
			},

			mockLimitsRepository: func() repository.LimitsRepositoryProvider {
				lrMock := limitsRepositoryMock.NewLimitsRepositoryProvider(s.T())

				lrMock.On(
					"CreateEntity",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
						) error {
							return nil
						},
					).
					Once()

				lrMock.On(
					"GetTransactionLimitsWithTotalSpendsByLimitType",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							limitType *model.LimitType,
						) []*repository.TransactionLimitWithTotalSpend {
							return nil
						},
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							limitType *model.LimitType,
						) error {
							return nil
						},
					).
					Once()

				return lrMock
			}(),

			mockTransactionLimitsUpdateRepository: func() repository.TransactionLimitsUpdateRepositoryProvider {
				tlurMock := limitsRepositoryMock.NewTransactionLimitsUpdateRepositoryProvider(s.T())
				tlurMock.On(
					"AddTransactionLimit",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							transactionLimit *model.TransactionLimit,
						) error {
							return nil
						},
					).
					Once()

				return tlurMock
			}(),

			mockMonoLithHttpInterface: func() client.MonolitClientProvider {
				lrMock := monolithHttpMock.NewMonolithHttpProvider(s.T())

				lrMock.On(
					"GetHistoricalTotalSpendsForLimitType",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							entity *model.Entity,
							limitType model.LimitType,
						) float64 {
							return 0
						},
						func(
							ctx context.Context,
							entity *model.Entity,
							limitType model.LimitType,
						) error {
							return nil
						},
					).
					Once()

				return lrMock
			}(),

			expectedLog: []string{
				fmt.Sprintf(literal.AddingTransactionLimit, entity.EntityType, entity.EntityId, model.Daily, 33.33),
			},
		},
		{
			name: "Adding New Limit That Exceeds Current Limit",

			transactionLimits: []*model.TransactionLimit{
				{
					LimitType:    model.Daily,
					Amount:       22.22,
					CurrencyCode: "SGD",
				},
			},

			// mockTimeInterface: func() timeInterface {
			// 	m := mocks.NewTimeInterface(s.T())
			// 	m.On("Now").
			// 		Return(
			// 			func() time.Time {
			// 				return time.Date(2001, 01, 01, 0, 0, 0, 0, loc)
			// 			},
			// 		)
			// 	return m
			// }(),

			mockLimitsRepository: func() repository.LimitsRepositoryProvider {
				lrMock := limitsRepositoryMock.NewLimitsRepositoryProvider(s.T())

				lrMock.On(
					"CreateEntity",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
						) error {
							return nil
						},
					).
					Once()

				lrMock.On(
					"GetTransactionLimitsWithTotalSpendsByLimitType",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							limitType *model.LimitType,
						) []*repository.TransactionLimitWithTotalSpend {
							return []*repository.TransactionLimitWithTotalSpend{
								{
									TotalSpend: 44.44,
									TransactionLimit: model.TransactionLimit{
										LimitType:    model.Daily,
										Amount:       33.33,
										CurrencyCode: "SGD",
									},
								},
							}
						},
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							limitType *model.LimitType,
						) error {
							return nil
						},
					).
					Once()

				return lrMock
			}(),

			mockTransactionLimitsUpdateRepository: func() repository.TransactionLimitsUpdateRepositoryProvider {
				tlurMock := limitsRepositoryMock.NewTransactionLimitsUpdateRepositoryProvider(s.T())

				tlurMock.On(
					"UpdateTransactionLimitAmount",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					float64(22.22),
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							transactionLimitId int64,
							amount float64,
						) error {
							return nil
						},
					).
					Once()

				return tlurMock
			}(),

			mockTimeInterface: func() transactionLimitsDateTimeHelper.TimeHelper {
				m := mocks.NewTimeHelper(s.T())

				m.On("Now").
					Return(
						func() time.Time {
							return time.Date(2001, 0o1, 0o1, 0, 0, 0, 0, loc)
						},
					).
					Once()

				return m
			}(),

			expectedLog: []string{
				fmt.Sprintf(literal.AddingTransactionLimit, entity.EntityType, entity.EntityId, model.Daily, 22.22),
				fmt.Sprintf(literal.TransactionLimitsAlreadyExistWillUpdate, model.Daily, 22.22),
			},

			out: []*service.UpdateTransactionLimitsResult{
				{
					PreviousLimit: 33.33,
					TotalSpends:   44.44,
					TransactionLimit: &model.TransactionLimit{
						LimitType:    model.Daily,
						CurrencyCode: "SGD",
						Amount:       22.22,
					},
					Message: fmt.Sprintf(literal.NewLimitAlreadyExceededLimit, model.Daily,
						time.Date(2001, 0o1, 0o2, 0, 0, 0, 0, loc).Format(literal.DateFormat)),
				},
			},
		},
		{
			name: "Error - Adding New Limit",

			transactionLimits: []*model.TransactionLimit{
				{
					LimitType:    model.Daily,
					Amount:       33.33,
					CurrencyCode: "SGD",
				},
			},

			mockLimitsRepository: func() repository.LimitsRepositoryProvider {

				lrMock := limitsRepositoryMock.NewLimitsRepositoryProvider(s.T())

				lrMock.On(
					"CreateEntity",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
						) error {
							return nil
						},
					).
					Once()

				lrMock.On(
					"GetTransactionLimitsWithTotalSpendsByLimitType",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							limitType *model.LimitType,
						) []*repository.TransactionLimitWithTotalSpend {
							return nil
						},
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
							limitType *model.LimitType,
						) error {
							return errors.New("error")
						},
					).
					Once()

				return lrMock

			}(),

			err: errors.New("error"),

			expectedLog: []string{
				fmt.Sprintf(literal.AddingTransactionLimit, entity.EntityType, entity.EntityId, model.Daily, 33.33),
			},
		},
		{
			name: "Error - on create entity",

			transactionLimits: []*model.TransactionLimit{
				{
					LimitType:    model.Daily,
					Amount:       33.33,
					CurrencyCode: "SGD",
				},
			},

			mockLimitsRepository: func() repository.LimitsRepositoryProvider {

				lrMock := limitsRepositoryMock.NewLimitsRepositoryProvider(s.T())

				lrMock.On(
					"CreateEntity",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						func(
							ctx context.Context,
							txn *sqlx.Tx,
							entity *model.Entity,
						) error {
							return errors.New("error")
						},
					).
					Once()

				return lrMock

			}(),

			err: errors.New("error"),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			var logBuffer bytes.Buffer
			logWriter := bufio.NewWriter(&logBuffer)
			log := log.New(logWriter, "", 0)
			// log, logBuffer, logWriter := support.NewTestLogger()
			log := support.NewLogger()
			log.Wr

			logWriter.Flush()

			assert.Equal(s.T(), test.out, msgOut)
			assert.Equal(s.T(), test.err, err)

			for _, expectedLog := range test.expectedLog {
				assert.Contains(s.T(), logBuffer.String(), expectedLog)
			}
		})
	}
}
