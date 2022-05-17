package exponent_retry

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 1) go test -v -cover -gcflags=all=-l -coverprofile=coverage.out
// 2) go tool cover -html=coverage.out

func TestDo(t *testing.T) {
	var errContentNum = -1
	var normalNil = -1
	var firstRetryNil = -1
	var secondRetryNil = -1
	var thridRetryNil = -1

	var ctxFirstRetryNil = -1
	var ctx context.Context

	ctx, _ = context.WithTimeout(context.Background(), time.Millisecond*1)

	var ctx2 context.Context
	var cancel2 context.CancelFunc
	ctx2, cancel2 = context.WithCancel(context.Background())

	type args struct {
		ctx context.Context
		cmd Cmd
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		errContent string
	}{
		{
			name: "test cancel inside",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					cancel2()
					if ctx2.Err() != nil {
						return ctx2.Err()
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "context canceled",
		},
		{
			name: "test deadline inside",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					if ctx.Err() != nil {
						return ctx.Err()
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "context deadline exceeded",
		},
		{
			name: "test deadline outside",
			args: args{
				ctx: ctx,
				cmd: func() error {
					ctxFirstRetryNil++

					switch ctxFirstRetryNil {
					case 0:
						return errors.New("normal execution")
					case 1:
						return nil
					case 2, 3:
						return errors.New("unimportance")
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "context deadline exceeded",
		},
		{
			name: "test normal nil",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					normalNil++

					switch normalNil {
					case 0:
						return nil
					case 1, 2, 3:
						return errors.New("unimportance")
					}

					return errors.New("not reached")
				},
			},
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test first retry nil",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					firstRetryNil++

					switch firstRetryNil {
					case 0:
						return errors.New("normal execution")
					case 1:
						return nil
					case 2, 3:
						return errors.New("unimportance")
					}

					return errors.New("not reached")
				},
			},
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test second retry nil",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					secondRetryNil++

					switch secondRetryNil {
					case 0:
						return errors.New("normal execution")
					case 1:
						return errors.New("first retry")
					case 2:
						return nil
					case 3:
						return errors.New("unimportance")
					}

					return errors.New("not reached")
				},
			},
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test third retry nil",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					thridRetryNil++

					switch thridRetryNil {
					case 0:
						return errors.New("normal execution")
					case 1:
						return errors.New("first retry")
					case 2:
						return errors.New("second retry")
					case 3:
						return nil
					}

					return errors.New("not reached")
				},
			},
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test err",
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					errContentNum++

					switch errContentNum {
					case 0, 1, 2:
						return errors.New("unimportance")
					case 3:
						return errors.New("third retry")
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "third retry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Do(tt.args.ctx, tt.args.cmd)

			assert.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("??? %v", err))

			if err != nil {
				assert.Equal(t, tt.errContent, err.Error())
			}
		})
	}
}

func TestDoWithReturn(t *testing.T) {
	var errContentNum = -1
	var normalNil = -1
	var firstRetryNil = -1
	var secondRetryNil = -1
	var thridRetryNil = -1

	var ctxFirstRetryNil = -1
	var ctx context.Context
	ctx, _ = context.WithTimeout(context.Background(), time.Millisecond*1)

	var ctx2 context.Context
	var cancel2 context.CancelFunc
	ctx2, cancel2 = context.WithCancel(context.Background())

	type args struct {
		ctx           context.Context
		cmdWithReturn CmdWithReturn
	}
	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		errContent string
	}{
		{
			name: "test cancel inside",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					cancel2()
					if ctx2.Err() != nil {
						return nil, ctx2.Err()
					}

					return 0, errors.New("not reached")
				},
			},
			want:       nil,
			wantErr:    true,
			errContent: "context canceled",
		},
		{
			name: "test deadline inside",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					if ctx.Err() != nil {
						return nil, ctx.Err()
					}

					return 0, errors.New("not reached")
				},
			},
			want:       nil,
			wantErr:    true,
			errContent: "context deadline exceeded",
		},
		{
			name: "test deadline outside",
			args: args{
				ctx: ctx,
				cmdWithReturn: func() (interface{}, error) {
					ctxFirstRetryNil++

					switch ctxFirstRetryNil {
					case 0:
						return ctxFirstRetryNil, errors.New("normal execution")
					case 1:
						return ctxFirstRetryNil, nil
					case 2, 3:
						return ctxFirstRetryNil, errors.New("unimportance")
					}

					return ctxFirstRetryNil, errors.New("not reached")
				},
			},
			want:       nil,
			wantErr:    true,
			errContent: "context deadline exceeded",
		},
		{
			name: "test normal nil",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					normalNil++

					switch normalNil {
					case 0:
						return normalNil, nil
					case 1, 2, 3:
						return normalNil, errors.New("unimportance")
					}

					return normalNil, errors.New("not reached")
				},
			},
			want:       0,
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test first retry nil",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					firstRetryNil++

					switch firstRetryNil {
					case 0:
						return firstRetryNil, errors.New("normal execution")
					case 1:
						return firstRetryNil, nil
					case 2, 3:
						return firstRetryNil, errors.New("unimportance")
					}

					return firstRetryNil, errors.New("not reached")
				},
			},
			want:       1,
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test second retry nil",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					secondRetryNil++

					switch secondRetryNil {
					case 0:
						return secondRetryNil, errors.New("normal execution")
					case 1:
						return secondRetryNil, errors.New("first retry")
					case 2:
						return secondRetryNil, nil
					case 3:
						return secondRetryNil, errors.New("unimportance")
					}

					return secondRetryNil, errors.New("not reached")
				},
			},
			want:       2,
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test third retry nil",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					thridRetryNil++

					switch thridRetryNil {
					case 0:
						return thridRetryNil, errors.New("normal execution")
					case 1:
						return thridRetryNil, errors.New("first retry")
					case 2:
						return thridRetryNil, errors.New("second retry")
					case 3:
						return thridRetryNil, nil
					}

					return thridRetryNil, errors.New("not reached")
				},
			},
			want:       3,
			wantErr:    false,
			errContent: "",
		},
		{
			name: "test err",
			args: args{
				ctx: context.Background(),
				cmdWithReturn: func() (interface{}, error) {
					errContentNum++

					switch errContentNum {
					case 0, 1, 2:
						return errContentNum, errors.New("unimportance")
					case 3:
						return errContentNum, errors.New("third retry")
					}

					return errContentNum, errors.New("not reached")
				},
			},
			want:       3,
			wantErr:    true,
			errContent: "third retry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoWithReturn(tt.args.ctx, tt.args.cmdWithReturn)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got, fmt.Sprintf("ff %v,%v", err, got))

			if err != nil {
				assert.Equal(t, tt.errContent, err.Error())
			}

		})
	}
}

func TestExponentRetry_Do(t *testing.T) {
	var deadlineOutsideN = -1
	var backoffN = -1
	var fnN = -1

	var ctx1 context.Context
	ctx1, _ = context.WithTimeout(context.Background(), time.Millisecond*88)

	var currentTime = time.Now()
	type fields struct {
		opt *Options
	}

	type args struct {
		ctx context.Context
		cmd Cmd
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		errContent string
	}{
		{
			name: "test deadline outside",
			fields: fields{
				&Options{
					MaxRetries:      1,
					MinRetryBackoff: time.Millisecond * 100,
					MaxRetryBackoff: time.Millisecond * 512,
				},
			},
			args: args{
				ctx: ctx1,
				cmd: func() error {
					deadlineOutsideN++

					switch deadlineOutsideN {
					case 0:
						return errors.New("unimportance")
					case 1:
						return nil
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "context deadline exceeded",
		},
		{
			name: "test backoff",
			fields: fields{
				&Options{
					MaxRetries:      11,
					MinRetryBackoff: time.Millisecond,
					MaxRetryBackoff: time.Millisecond * 1024,
				},
			},
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					backoffN++

					switch backoffN {
					//   1  2  4  8  16 32 64 128 256 512 1024
					//   0  1  2  3  4  5  6  7   8   9   10  11
					case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
						return errors.New("unimportance")
					case 11: // 1024
						if time.Now().Sub(currentTime) > (time.Millisecond * 2047) {
							return errors.New("time ok")
						}

						return errors.New("time not ok" + fmt.Sprintf("%v", time.Now().Sub(currentTime)))
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "time ok",
		},
		{
			name: "test fn",
			fields: fields{
				&Options{
					MaxRetries:      5,
					MinRetryBackoff: time.Millisecond,
					MaxRetryBackoff: time.Millisecond * 1024,
					Fn: func(err error) bool {
						if err.Error() == "pass" {
							return false
						}

						return true
					},
				},
			},
			args: args{
				ctx: context.Background(),
				cmd: func() error {
					fnN++

					switch fnN {
					case 0:
						return errors.New("normal execution")
					case 1:
						return errors.New("pass")
					case 2:
						return errors.New("second retry")
					}

					return errors.New("not reached")
				},
			},
			wantErr:    true,
			errContent: "pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := &ExponentRetry{
				opt: tt.fields.opt,
			}
			err := er.Do(tt.args.ctx, tt.args.cmd)

			assert.Equal(t, tt.wantErr, err != nil)

			if err != nil {
				assert.Equal(t, tt.errContent, err.Error())
			}
		})
	}
}

func TestExponentRetry_DoWithReturn(t *testing.T) {
	type fields struct {
		opt *Options
	}
	type args struct {
		ctx           context.Context
		cmdWithReturn CmdWithReturn
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := &ExponentRetry{
				opt: tt.fields.opt,
			}
			got, err := er.DoWithReturn(tt.args.ctx, tt.args.cmdWithReturn)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoWithReturn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoWithReturn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts *Options
	}
	tests := []struct {
		name string
		args args
		want *ExponentRetry
	}{
		{
			name: "test options",
			args: args{
				&Options{},
			},
			want: &ExponentRetry{
				&Options{
					MaxRetries:      3,
					MinRetryBackoff: time.Millisecond * 8,
					MaxRetryBackoff: time.Millisecond * 512,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.opts)

			assert.Equal(t, tt.want, got)
			assert.NotSame(t, tt.want, got)
		})
	}
}
