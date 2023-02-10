package entities

type SlidingWindowResponseChannel struct {
	ServiceId     int
	ShouldForward bool
	ErrorMessage  string
}
