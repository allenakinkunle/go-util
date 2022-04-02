package identifier

// mockNanoGenerator Mock of IIdentifier
type mockIdentifier struct{}

func NewMockIdentifier() mockIdentifier {
	return mockIdentifier{}
}

func (m mockIdentifier) NewNanoID() string {
	return "newnanoid"
}

func (m mockIdentifier) NewUUIDv5(name string) string {
	return "newuuidv5"
}
