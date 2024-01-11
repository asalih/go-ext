package disklayout

import "testing"

func Test_getExtType(t *testing.T) {
	type args struct {
		fCompat   uint32
		fInCompat uint32
	}
	tests := []struct {
		name string
		args args
		want ExtType
	}{
		{name: "ext4", args: args{fCompat: 0x3c, fInCompat: 0x2c6}, want: Ext4},
		{name: "ext3", args: args{fCompat: 0x3c, fInCompat: 0x6}, want: Ext3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExtType(tt.args.fCompat, tt.args.fInCompat); got != tt.want {
				t.Errorf("getExtType() = %v, want %v", got, tt.want)
			}
		})
	}
}
