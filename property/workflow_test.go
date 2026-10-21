package property

import "testing"

func TestVerifyLoginDecision(t *testing.T) {
	tests := []struct { name string; verified bool; tenant string; want bool }{
		{"verified tenant", true, "tenant-17", true}, {"rejected code", false, "tenant-17", false}, {"missing tenant", true, "", false},
	}
	for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { got, err := VerifyLogin(tt.verified, tt.tenant); if (err == nil) != tt.want || got.Verified != tt.want { t.Fatalf("got decision=%+v err=%v", got, err) } }) }
}
