package interfacepollution

// Consumer defines interface with only what it needs
// This is where interface is "discovered"

// UserCreator interface defined by consumer
type UserCreator interface {
	CreateUser(name string) error
}

// UserReader interface for different consumer needs
type UserReader interface {
	GetUser(id int) (string, error)
}

// Handler only depends on what it needs
type UserHandler struct {
	creator UserCreator
	reader  UserReader
}

func NewUserHandler(creator UserCreator, reader UserReader) *UserHandler {
	return &UserHandler{
		creator: creator,
		reader:  reader,
	}
}

func (h *UserHandler) HandleCreate(name string) error {
	return h.creator.CreateUser(name)
}

func (h *UserHandler) HandleGet(id int) (string, error) {
	return h.reader.GetUser(id)
}

// UserService automatically satisfies both interfaces
// No explicit "implements" needed

// Benefits:
// 1. Interface defined where it's used
// 2. Only required methods exposed
// 3. Easy to mock for testing
// 4. UserService doesn't know about Handler
