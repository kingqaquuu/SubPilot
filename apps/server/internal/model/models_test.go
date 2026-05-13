package model

import (
	"reflect"
	"testing"
)

func TestAllModels(t *testing.T) {
	models := All()
	if len(models) != 5 {
		t.Fatalf("expected 5 models, got %d", len(models))
	}

	expected := []reflect.Type{
		reflect.TypeOf(&User{}),
		reflect.TypeOf(&Category{}),
		reflect.TypeOf(&Subscription{}),
		reflect.TypeOf(&Reminder{}),
		reflect.TypeOf(&PaymentRecord{}),
	}
	for i := range expected {
		if reflect.TypeOf(models[i]) != expected[i] {
			t.Fatalf("unexpected model at index %d: %T", i, models[i])
		}
	}
}
