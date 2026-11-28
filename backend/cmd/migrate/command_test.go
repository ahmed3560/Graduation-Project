package main

import "testing"

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    command
		wantErr bool
	}{
		{"up", []string{"up"}, command{name: "up"}, false},
		{"version", []string{"version"}, command{name: "version"}, false},
		{"down 1", []string{"down", "1"}, command{name: "down", steps: 1}, false},
		{"down 12", []string{"down", "12"}, command{name: "down", steps: 12}, false},

		{"no arguments", nil, command{}, true},
		{"unknown command", []string{"sideways"}, command{}, true},
		{"up with extra argument", []string{"up", "1"}, command{}, true},
		{"version with extra argument", []string{"version", "x"}, command{}, true},
		{"down without a number", []string{"down"}, command{}, true},
		{"down zero", []string{"down", "0"}, command{}, true},
		{"down negative", []string{"down", "-1"}, command{}, true},
		{"down not a number", []string{"down", "all"}, command{}, true},
		{"down two numbers", []string{"down", "1", "2"}, command{}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseCommand(tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("command = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestCheckAllowed(t *testing.T) {
	down := command{name: "down", steps: 1}
	up := command{name: "up"}

	if err := down.checkAllowed("production"); err == nil {
		t.Error("down must be refused in production")
	}
	for _, env := range []string{"development", "test"} {
		if err := down.checkAllowed(env); err != nil {
			t.Errorf("down in %s: unexpected error %v", env, err)
		}
	}
	if err := up.checkAllowed("production"); err != nil {
		t.Errorf("up in production: unexpected error %v", err)
	}
}
