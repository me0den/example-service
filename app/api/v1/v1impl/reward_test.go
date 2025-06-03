package v1impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"

	v1 "github.com/me0den/example-service/app/api/v1"
	"github.com/me0den/example-service/app/api/v1/transport/routes"
	"github.com/me0den/example-service/app/api/v1/v1impl/mock"
	"github.com/me0den/example-service/domain/entity"
)

func TestRewardService_CreateReward(t *testing.T) {
	type args struct {
		ctx       echo.Context
		req       *v1.CreateRewardRequest
		isBadBind bool
	}

	type listUserElosArgs struct {
		teams []*entity.Team
	}

	type listUserElosWant struct {
		err      error
		userElos []*entity.UserElo
	}

	type calculateEloArgs struct {
		userElos  []*entity.UserElo
		winnerIdx int
	}

	type calculateEloWant struct {
		newUserElos []*entity.UserElo
	}

	type BatchUpdateEloArgs struct {
		newUserElos []*entity.UserElo
	}

	type BatchUpdateEloWant struct {
		err error
	}

	defaultElo := 1000
	tests := []struct {
		name               string
		args               args
		err                error
		want               *v1.CreateRewardResponse
		wantErr            bool
		listUserElosArgs   *listUserElosArgs
		listUserElosWant   *listUserElosWant
		calculateEloArgs   *calculateEloArgs
		calculateEloWant   *calculateEloWant
		BatchUpdateEloArgs *BatchUpdateEloArgs
		BatchUpdateEloWant *BatchUpdateEloWant
	}{
		{
			name: "Validator request form: required validation fail",
			args: args{
				req: &v1.CreateRewardRequest{},
			},
			want:    nil,
			err:     echo.NewHTTPError(http.StatusBadRequest, []string{"winner is required", "teams is required"}),
			wantErr: true,
		},
		{
			name: "Validator request form: teams is not equals to 2.",
			args: args{
				req: &v1.CreateRewardRequest{
					Teams: []*entity.Team{
						{
							ID:    "team_1",
							Owner: "user_1",
						},
					},
					Winner: "user_1",
				},
			},
			want:    nil,
			err:     echo.NewHTTPError(http.StatusBadRequest, []string{"teams must be equals to 2"}),
			wantErr: true,
		},
		{
			name: "Cannot bind http request body.",
			args: args{
				req:       &v1.CreateRewardRequest{},
				isBadBind: true,
			},
			want:    nil,
			err:     echo.NewHTTPError(http.StatusBadRequest, "code=415, message=Unsupported Media Type"),
			wantErr: true,
		},
		{
			name: "Successful create reward: user_1 is winner",
			args: args{
				req: &v1.CreateRewardRequest{
					Teams: []*entity.Team{
						{
							ID:    "team_1",
							Owner: "user_1",
						},
						{
							ID:    "team_2",
							Owner: "user_2",
						},
					},
					Winner: "user_1",
				},
			},
			want: &v1.CreateRewardResponse{
				Items: []*v1.Reward{
					{
						UserID: "user_1",
						OldElo: defaultElo,
						NewElo: defaultElo + 10,
					},
					{
						UserID: "user_2",
						OldElo: defaultElo,
						NewElo: defaultElo - 10,
					},
				},
			},
			err:     nil,
			wantErr: false,
			listUserElosArgs: &listUserElosArgs{
				teams: []*entity.Team{
					{
						ID:    "team_1",
						Owner: "user_1",
					},
					{
						ID:    "team_2",
						Owner: "user_2",
					},
				},
			},
			listUserElosWant: &listUserElosWant{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
			},
			calculateEloArgs: &calculateEloArgs{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
				winnerIdx: 1,
			},
			calculateEloWant: &calculateEloWant{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo + 10,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo - 10,
					},
				},
			},
			BatchUpdateEloArgs: &BatchUpdateEloArgs{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo + 10,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo - 10,
					},
				},
			},
			BatchUpdateEloWant: &BatchUpdateEloWant{
				err: nil,
			},
		},
		{
			name: "Successful create reward: user_2 is winner",
			args: args{
				req: &v1.CreateRewardRequest{
					Teams: []*entity.Team{
						{
							ID:    "team_1",
							Owner: "user_1",
						},
						{
							ID:    "team_2",
							Owner: "user_2",
						},
					},
					Winner: "user_2",
				},
			},
			want: &v1.CreateRewardResponse{
				Items: []*v1.Reward{
					{
						UserID: "user_1",
						OldElo: defaultElo,
						NewElo: defaultElo - 10,
					},
					{
						UserID: "user_2",
						OldElo: defaultElo,
						NewElo: defaultElo + 10,
					},
				},
			},
			err:     nil,
			wantErr: false,
			listUserElosArgs: &listUserElosArgs{
				teams: []*entity.Team{
					{
						ID:    "team_1",
						Owner: "user_1",
					},
					{
						ID:    "team_2",
						Owner: "user_2",
					},
				},
			},
			listUserElosWant: &listUserElosWant{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
			},
			calculateEloArgs: &calculateEloArgs{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
				winnerIdx: 2,
			},
			calculateEloWant: &calculateEloWant{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo - 10,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo + 10,
					},
				},
			},
			BatchUpdateEloArgs: &BatchUpdateEloArgs{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo - 10,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo + 10,
					},
				},
			},
			BatchUpdateEloWant: &BatchUpdateEloWant{
				err: nil,
			},
		},
		{
			name: "Successful create reward: draw",
			args: args{
				req: &v1.CreateRewardRequest{
					Teams: []*entity.Team{
						{
							ID:    "team_1",
							Owner: "user_1",
						},
						{
							ID:    "team_2",
							Owner: "user_2",
						},
					},
					Winner: "draw",
				},
			},
			want: &v1.CreateRewardResponse{
				Items: []*v1.Reward{
					{
						UserID: "user_1",
						OldElo: defaultElo,
						NewElo: defaultElo + 5,
					},
					{
						UserID: "user_2",
						OldElo: defaultElo,
						NewElo: defaultElo + 5,
					},
				},
			},
			err:     nil,
			wantErr: false,
			listUserElosArgs: &listUserElosArgs{
				teams: []*entity.Team{
					{
						ID:    "team_1",
						Owner: "user_1",
					},
					{
						ID:    "team_2",
						Owner: "user_2",
					},
				},
			},
			listUserElosWant: &listUserElosWant{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
			},
			calculateEloArgs: &calculateEloArgs{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
				winnerIdx: 2,
			},
			calculateEloWant: &calculateEloWant{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo + 5,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo + 5,
					},
				},
			},
			BatchUpdateEloArgs: &BatchUpdateEloArgs{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo + 5,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo + 5,
					},
				},
			},
			BatchUpdateEloWant: &BatchUpdateEloWant{
				err: nil,
			},
		},
		{
			name: "Failed to get user elo",
			args: args{
				req: &v1.CreateRewardRequest{
					Teams: []*entity.Team{
						{
							ID:    "team_1",
							Owner: "user_1",
						},
						{
							ID:    "team_2",
							Owner: "user_2",
						},
					},
					Winner: "draw",
				},
			},
			want:    nil,
			err:     echo.ErrInternalServerError,
			wantErr: true,
			listUserElosArgs: &listUserElosArgs{
				teams: []*entity.Team{
					{
						ID:    "team_1",
						Owner: "user_1",
					},
					{
						ID:    "team_2",
						Owner: "user_2",
					},
				},
			},
			listUserElosWant: &listUserElosWant{
				userElos: nil,
				err:      echo.ErrInternalServerError,
			},
		},
		{
			name: "Failed to batch update elo",
			args: args{
				req: &v1.CreateRewardRequest{
					Teams: []*entity.Team{
						{
							ID:    "team_1",
							Owner: "user_1",
						},
						{
							ID:    "team_2",
							Owner: "user_2",
						},
					},
					Winner: "user_1",
				},
			},
			want:    nil,
			err:     echo.ErrInternalServerError,
			wantErr: true,
			listUserElosArgs: &listUserElosArgs{
				teams: []*entity.Team{
					{
						ID:    "team_1",
						Owner: "user_1",
					},
					{
						ID:    "team_2",
						Owner: "user_2",
					},
				},
			},
			listUserElosWant: &listUserElosWant{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
			},
			calculateEloArgs: &calculateEloArgs{
				userElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo,
					},
				},
				winnerIdx: 2,
			},
			calculateEloWant: &calculateEloWant{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo + 10,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo - 10,
					},
				},
			},
			BatchUpdateEloArgs: &BatchUpdateEloArgs{
				newUserElos: []*entity.UserElo{
					{
						UserID: "user_1",
						Elo:    defaultElo + 10,
					},
					{
						UserID: "user_2",
						Elo:    defaultElo - 10,
					},
				},
			},
			BatchUpdateEloWant: &BatchUpdateEloWant{
				err: echo.ErrInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		e := echo.New()
		rec := httptest.NewRecorder()
		ctx := context.Background()
		if tt.args.req != nil {
			marshalled, _ := json.Marshal(tt.args.req)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshalled))
			if !tt.args.isBadBind {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			tt.args.ctx = e.NewContext(req, rec)
		}

		validate := validator.New()
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "-" {
				return ""
			}
			return name
		})

		e.Validator = &routes.Validator{Validator: validate}

		t.Run(tt.name, func(t *testing.T) {
			redisRepo := &mock.RedisRepo{}
			svcMock := &mock.RewardService{}
			svc := &RewardService{
				redisRepo: redisRepo,
			}

			if tt.listUserElosArgs != nil && tt.listUserElosWant != nil {
				for idx, team := range tt.listUserElosArgs.teams {
					userElo := entity.NewUserDefaultElo(team.Owner)
					redisRepo.On("GetUserElo", ctx, team.Owner).Return(userElo, tt.listUserElosWant.err)
					if tt.listUserElosWant.err == nil {
						tt.listUserElosWant.userElos[idx] = userElo
					}
				}
			}

			if tt.calculateEloArgs != nil && tt.calculateEloWant != nil {
				svcMock.On("calculateElo", tt.args.ctx, tt.calculateEloArgs.userElos, tt.calculateEloArgs.winnerIdx).
					Return(tt.calculateEloWant.newUserElos)
			}

			if tt.BatchUpdateEloArgs != nil && tt.BatchUpdateEloWant != nil {
				redisRepo.On("BatchUpdateElo", ctx, tt.BatchUpdateEloArgs.newUserElos).
					Return(tt.BatchUpdateEloWant.err)
			}

			beforeTime := time.Now().Unix()
			err := svc.CreateReward(tt.args.ctx)
			afterTime := time.Now().Unix()
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.want != nil {
				var resp v1.CreateRewardResponse
				err = json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				wantMarshal, err := json.Marshal(tt.want)
				assert.NoError(t, err)
				// Verify timestamp is within reasonable range
				assert.True(t, resp.Items[0].UpdatedAt >= beforeTime)
				assert.True(t, resp.Items[0].UpdatedAt <= afterTime)

				// Remove time field in response to compare
				for _, item := range resp.Items {
					item.UpdatedAt = 0
				}

				respMarshal, err := json.Marshal(resp)
				assert.NoError(t, err)
				assert.Equal(t, string(wantMarshal), string(respMarshal))
			}
		})
	}
}

func TestRewardService_calculateElo(t *testing.T) {
	type args struct {
		ctx       context.Context
		userElos  []*entity.UserElo
		winnerIdx int
	}
	tests := []struct {
		name string
		args args
		want []*entity.UserElo
	}{
		{
			name: "Winner index 0 - tie case",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: 1000},
					{UserID: "user_2", Elo: 1200},
				},
				winnerIdx: 0,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 1005},
				{UserID: "user_2", Elo: 1205},
			},
		},
		{
			name: "Winner index 1 - player 0 wins",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: 1000},
					{UserID: "user_2", Elo: 1200},
				},
				winnerIdx: 1,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 1010},
				{UserID: "user_2", Elo: 1190},
			},
		},
		{
			name: "Winner index 2 - player 1 wins",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: 1000},
					{UserID: "user_2", Elo: 1200},
				},
				winnerIdx: 2,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 990},
				{UserID: "user_2", Elo: 1210},
			},
		},
		{
			name: "Zero Elo values",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: 0},
					{UserID: "user_2", Elo: 0},
				},
				winnerIdx: 1,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 10},
				{UserID: "user_2", Elo: -10},
			},
		},
		{
			name: "Negative Elo values",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: -50},
					{UserID: "user_2", Elo: -30},
				},
				winnerIdx: 2,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: -60},
				{UserID: "user_2", Elo: -20},
			},
		},
		{
			name: "Large Elo values",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: 2500},
					{UserID: "user_2", Elo: 2800},
				},
				winnerIdx: 0,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 2505},
				{UserID: "user_2", Elo: 2805},
			},
		},
		{
			name: "Invalid negative index",
			args: args{
				userElos: []*entity.UserElo{
					{UserID: "user_1", Elo: 1000},
					{UserID: "user_2", Elo: 1200},
				},
				winnerIdx: -1,
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 1000},
				{UserID: "user_2", Elo: 1200},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &RewardService{}
			res := svc.calculateElo(tt.args.ctx, tt.args.userElos, tt.args.winnerIdx)
			if tt.want != nil {
				assert.Equal(t, tt.want, res)
			}
		})
	}
}

func TestRewardService_listUserElos(t *testing.T) {
	type args struct {
		ctx   context.Context
		teams []*entity.Team
	}
	tests := []struct {
		name       string
		args       args
		err        error
		want       []*entity.UserElo
		wantErr    bool
		setupMocks func(repo *mock.RedisRepo)
	}{
		{
			name: "successful retrieval of user elos",
			args: args{
				teams: []*entity.Team{
					{Owner: "user_1"},
					{Owner: "user_2"},
				},
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 1000},
				{UserID: "user_2", Elo: 1000},
			},
			err:     nil,
			wantErr: false,
			setupMocks: func(mockRepo *mock.RedisRepo) {
				mockRepo.On("GetUserElo", tmock.Anything, "user_1").Return(&entity.UserElo{UserID: "user_1", Elo: 1000}, nil)
				mockRepo.On("GetUserElo", tmock.Anything, "user_2").Return(&entity.UserElo{UserID: "user_2", Elo: 1000}, nil)
			},
		},
		{
			name: "single team",
			args: args{
				teams: []*entity.Team{
					{Owner: "user_1"},
				},
			},
			setupMocks: func(mockRepo *mock.RedisRepo) {
				mockRepo.On("GetUserElo", tmock.Anything, "user_1").Return(&entity.UserElo{UserID: "user_1", Elo: 1000}, nil)
			},
			want: []*entity.UserElo{
				{UserID: "user_1", Elo: 1000},
			},
			err:     nil,
			wantErr: false,
		},
		{
			name: "empty teams slice",
			args: args{
				teams: []*entity.Team{},
			},
			want:       nil,
			err:        nil,
			wantErr:    false,
			setupMocks: func(mockRepo *mock.RedisRepo) {},
		},
		{
			name: "nil teams slice",
			args: args{
				teams: nil,
			},
			want:       nil,
			err:        nil,
			wantErr:    false,
			setupMocks: func(mockRepo *mock.RedisRepo) {},
		},
		{
			name: "error from redis repository",
			args: args{
				teams: []*entity.Team{
					{Owner: "user_1"},
					{Owner: "user_2"},
				},
			},
			setupMocks: func(mockRepo *mock.RedisRepo) {
				mockRepo.On("GetUserElo", tmock.Anything, "user_1").Return(&entity.UserElo{UserID: "user_1", Elo: 1000}, nil)
				mockRepo.On("GetUserElo", tmock.Anything, "user_2").Return(nil, errors.New("redis connection failed"))
			},
			want:    nil,
			err:     errors.New("redis connection failed"),
			wantErr: true,
		},
		{
			name: "error on first user",
			args: args{
				teams: []*entity.Team{
					{Owner: "user_1"},
					{Owner: "user_2"},
				},
			},
			setupMocks: func(mockRepo *mock.RedisRepo) {
				mockRepo.On("GetUserElo", tmock.Anything, "user_1").Return(nil, errors.New("user not found"))
			},
			want:    nil,
			err:     errors.New("user not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			redisRepo := &mock.RedisRepo{}
			tt.setupMocks(redisRepo)
			svc := &RewardService{
				redisRepo: redisRepo,
			}
			rec, err := svc.listUserElos(ctx, tt.args.teams)
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want, rec)
			}
		})
	}
}
