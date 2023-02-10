package usecases_test

import (
	"context"
	"testing"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBTCUC_GetUserBalance(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.UserBalance
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Successfully get User balance": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := &rpc.UserBalance{}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(want, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.GetUserBalance(tt.args.ctx, tt.args.userID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
