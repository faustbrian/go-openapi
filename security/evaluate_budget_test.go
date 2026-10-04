package security

import "testing"

func TestConsumeBudgetAdmission(t *testing.T) {
	remaining := 5
	if !consumeBudget(&remaining, 3) || remaining != 2 {
		t.Fatal("ordinary occurrence did not consume its exact budget")
	}
	if consumeBudget(&remaining, 3) || remaining != 2 {
		t.Fatal("one-over occurrence changed the remaining budget")
	}
	if !consumeBudget(&remaining, 2) || remaining != 0 {
		t.Fatal("inclusive final occurrence was refused")
	}
	if !consumeBudget(&remaining, 0) || remaining != 0 {
		t.Fatal("empty occurrence changed an exhausted budget")
	}
	if consumeBudget(&remaining, 1) || remaining != 0 {
		t.Fatal("exhausted budget admitted a further occurrence")
	}
}
