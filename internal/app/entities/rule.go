package entities

type ServiceEntity struct {
	ServiceId int  `dynamodbav:"service_id"`
	Rule      Rule `dynamodbav:"rule"`
}

type Rule struct {
	Field                   string
	Rate                    Rate
	RequestRejectionMessage string
}

func (r Rule) GetField() string {
	return r.Field
}
