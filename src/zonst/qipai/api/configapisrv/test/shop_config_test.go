package test

import (
	"testing"
	"zonst/qipai/api/configapisrv/service/some_game_config_copy"
)

func Test_shopConfig_Copy(t *testing.T) {
	type args struct {
		param some_game_config_copy.Param
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: 复制道具系统.
		{
			"测试1",
			args{some_game_config_copy.Param{
				API:        "http://123.206.215.185:9960",
				JwtToken:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNTc0LCJ1c2VyX25hbWUiOiLmlZbnh5XlhbUiLCJncm91cF9pZCI6MiwiaXNfc3VwZXJ1c2VyIjpmYWxzZSwiZXhwIjoxNjEyNzYzMjcxfQ.sX0QE40129_wr4cNmunJmxCofgkO9UupQDiVpxvG1PA",
				SrcGameID:  67,
				DestGameID: 88,
			}},
			false,
		},
		{
			"测试2",
			args{some_game_config_copy.Param{
				API:        "http://123.206.215.185:9960",
				JwtToken:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNTc0LCJ1c2VyX25hbWUiOiLmlZbnh5XlhbUiLCJncm91cF9pZCI6MiwiaXNfc3VwZXJ1c2VyIjpmYWxzZSwiZXhwIjoxNjEyNzYzMjcxfQ.sX0QE40129_wr4cNmunJmxCofgkO9UupQDiVpxvG1PA",
				SrcGameID:  6767,
				DestGameID: 88,
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := some_game_config_copy.ShopConfig{}
			if err := s.Copy(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
