package errs

import (
	"errors"
	"testing"
)

func TestRpcErrCheck(t *testing.T) {
	type args struct {
		err      error
		checkErr CodeError
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "aaa",
			args: args{
				err:      errors.New("aaa"),
				checkErr: ErrDiamondNotEnough,
			},
			want: true,
		},
		{
			name: "bbb",
			args: args{
				err:      ErrDiamondNotEnough,
				checkErr: ErrDiamondNotEnough,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RpcErrCheck(tt.args.err, tt.args.checkErr); got != tt.want {
				t.Errorf("RpcErrCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
