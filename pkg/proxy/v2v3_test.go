package proxy

import "testing"

func Test_formatKey(t *testing.T) {
	type args struct {
		key []byte
		val []byte
	}
	tests := []struct {
		name       string
		args       args
		wantNewKey string
		wantNewVal string
		wantErr    bool
	}{
		{
			name: "std upsync format",
			args: args{
				key: []byte("/upstream/somerpc/127.0.0.1:8081"),
				val: []byte{},
			},
			wantNewKey: "/upstream/somerpc/127.0.0.1:8081",
			wantNewVal: "",
			wantErr:    false,
		},
		{
			name: "go-zero style",
			args: args{
				key: []byte("/service/hellorpc.rpc/2131421513123512"),
				val: []byte("10.254.78.4:8081"),
			},
			wantNewKey: "/service/hellorpc.rpc/10.254.78.4:8081",
			wantNewVal: "",
			wantErr:    false,
		},
		{
			name: "url style",
			args: args{
				key: []byte("/service/hellorpc.rpc/2131421513123512"),
				val: []byte("//127.0.0.1:8080/?weight=1&max_fails=2&fail_timeout=10"),
			},
			wantNewKey: "/service/hellorpc.rpc/127.0.0.1:8080",
			wantNewVal: "{\"weight\":1,\"max_fails\":2,\"fail_timeout\":10}",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewKey, gotNewVal, err := formatKey(tt.args.key, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewKey != tt.wantNewKey {
				t.Errorf("formatKey() gotNewKey = %v, want %v", gotNewKey, tt.wantNewKey)
			}
			if gotNewVal != tt.wantNewVal {
				t.Errorf("formatKey() gotNewVal = %v, want %v", gotNewVal, tt.wantNewVal)
			}
		})
	}
}
