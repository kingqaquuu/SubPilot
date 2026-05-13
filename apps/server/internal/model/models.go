package model

func All() []interface{} {
	return []interface{}{
		&User{},
		&Category{},
		&Subscription{},
		&Reminder{},
		&PaymentRecord{},
	}
}
