package foobarbaz

//go:generate go run github.com/golang/mock/mockgen -source foo_service.go -destination mock/foo_service_mock.go -package foobarbaz_mock

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/model"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

// FooService is the service interface for Foo entities.
type FooService interface {
	Create(requestFormat FooRequestFormat, userID uuid.UUID) (foo Foo, err error)
	ResolveByID(id uuid.UUID, withItems bool) (foo Foo, err error)
	SoftDelete(id uuid.UUID, userID uuid.UUID) (foo Foo, err error)
	Update(id uuid.UUID, requestFormat FooRequestFormat, userID uuid.UUID) (foo Foo, err error)
}

// FooServiceImpl is the service implementation for Foo entities.
type FooServiceImpl struct {
	FooRepository FooRepository
	Producer      producer.Producer
	Config        *configs.Config
}

// ProvideFooServiceImpl is the provider for this service.
func ProvideFooServiceImpl(fooRepository FooRepository, producer producer.Producer, config *configs.Config) *FooServiceImpl {
	s := new(FooServiceImpl)
	s.FooRepository = fooRepository
	s.Config = config
	s.Producer = producer

	return s
}

// Create creates a new Foo.
func (s *FooServiceImpl) Create(requestFormat FooRequestFormat, userID uuid.UUID) (foo Foo, err error) {
	foo, err = foo.NewFromRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}

	if err != nil {
		return foo, failure.BadRequest(err)
	}

	err = s.FooRepository.Create(foo)

	if err != nil {
		return
	}

	if s.Config.Event.Producer.SNS.Topics.FooCreated.Enabled {
		e := model.NewEvent(FooBarBazEventType, requestFormat)
		s.Producer.Publish(model.PublishRequest{
			Event: e,
			Topic: s.Config.Event.Producer.SNS.Topics.FooCreated.ARN,
		})
	}

	return
}

// ResolveByID resolves a Foo by its ID.
func (s *FooServiceImpl) ResolveByID(id uuid.UUID, withItems bool) (foo Foo, err error) {
	foo, err = s.FooRepository.ResolveByID(id)

	if foo.IsDeleted() {
		return foo, failure.NotFound("foo")
	}

	if withItems {
		items, err := s.FooRepository.ResolveItemsByFooIDs([]uuid.UUID{foo.ID})
		if err != nil {
			return foo, err
		}

		foo.AttachItems(items)
	}

	return
}

// SoftDelete marks a Foo as deleted by setting its `deleted` and `deletedBy` properties.
func (s *FooServiceImpl) SoftDelete(id uuid.UUID, userID uuid.UUID) (foo Foo, err error) {
	foo, err = s.FooRepository.ResolveByID(id)
	if err != nil {
		return
	}

	// need to get the items so they don't get deleted
	items, err := s.FooRepository.ResolveItemsByFooIDs([]uuid.UUID{foo.ID})
	if err != nil {
		return foo, err
	}

	foo.AttachItems(items)

	err = foo.SoftDelete(userID)
	if err != nil {
		return
	}

	err = s.FooRepository.Update(foo)
	return
}

// Update updates a Foo.
func (s *FooServiceImpl) Update(id uuid.UUID, requestFormat FooRequestFormat, userID uuid.UUID) (foo Foo, err error) {
	foo, err = s.FooRepository.ResolveByID(id)
	if err != nil {
		return
	}

	err = foo.Update(requestFormat, userID)
	if err != nil {
		return
	}

	err = s.FooRepository.Update(foo)
	return
}
