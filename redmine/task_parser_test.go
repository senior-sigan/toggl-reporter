package redmine

import "testing"

func TestParseTaskText(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{name: "Should find task", args: args{description: "Task #42 Some work"}, want: 42, want1: true},
		{name: "Should find task", args: args{description: "Story #42 Some work"}, want: 42, want1: true},
		{name: "Should find task", args: args{description: "Bug #42 Some work"}, want: 42, want1: true},
		{name: "Should find task", args: args{description: "38: Some work"}, want: 38, want1: true},
		{name: "Should find task", args: args{description: "38: Some work"}, want: 38, want1: true},
		{name: "Should find task", args: args{description: "#42 Some work"}, want: 0, want1: false},
		{name: "Should not find task", args: args{description: "Some task with 42 number"}, want: 0, want1: false},
		{name: "Should not find task", args: args{description: "Some task with 42: number"}, want: 0, want1: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindTaskId(tt.args.description)
			if got != tt.want {
				t.Errorf("FindTaskId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindTaskId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
