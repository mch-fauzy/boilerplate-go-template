package foobarbaz_test

import (
	"testing"
	"time"

	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	foobarbaz_mock "github.com/evermos/boilerplate-go/internal/domain/foobarbaz/mock"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

func getRandomUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func uuidFromString(s string) uuid.UUID {
	id, _ := uuid.FromString(s)
	return id
}

func TestFooService(t *testing.T) {

	t.Run("resolveByID", func(t *testing.T) {
		tests := []struct {
			name        string
			entityID    uuid.UUID
			setupMock   func(*foobarbaz_mock.MockFooRepository, uuid.UUID, foobarbaz.Foo, []foobarbaz.FooItem, error)
			returns     *foobarbaz.Foo
			returnItems *[]foobarbaz.FooItem
			err         error
		}{
			{
				name:     "default",
				entityID: uuidFromString("4e80c5bf-b79b-4c90-8f91-82647f439e55"),
				setupMock: func(mockRepo *foobarbaz_mock.MockFooRepository, id uuid.UUID, ent foobarbaz.Foo, entItems []foobarbaz.FooItem, err error) {
					mockRepo.EXPECT().ResolveByID(id).Return(ent, err)
					mockRepo.EXPECT().ResolveItemsByFooIDs([]uuid.UUID{id}).Return(entItems, err)
				},
				returns: &foobarbaz.Foo{
					ID:            uuidFromString("4e80c5bf-b79b-4c90-8f91-82647f439e55"),
					Name:          "The First Foo",
					TotalQuantity: int64(5),
					TotalPrice:    float64(65000),
					TotalDiscount: float64(3900),
					ShippingFee:   float64(15000),
					GrandTotal:    float64(76100),
					Status:        foobarbaz.FooStatusNew,
					Created:       time.Now(),
					CreatedBy:     getRandomUUID(),
					Updated:       null.TimeFrom(time.Now()),
					UpdatedBy:     nuuid.From(getRandomUUID()),
				},
				returnItems: &[]foobarbaz.FooItem{
					{
						ID:          uuidFromString("7e94b76a-0fc7-4422-bdb0-0caa2f80e43f"),
						FooID:       uuidFromString("4e80c5bf-b79b-4c90-8f91-82647f439e55"),
						SKU:         "SKU-00001",
						ProductName: "Product Name 1",
						Quantity:    int64(2),
						UnitPrice:   float64(10000),
						TotalPrice:  float64(20000),
						Discount:    float64(1200),
						GrandTotal:  float64(18800),
					},
					{
						ID:          uuidFromString("c43ce49f-c689-4f06-9f58-7dec2952beeb"),
						FooID:       uuidFromString("4e80c5bf-b79b-4c90-8f91-82647f439e55"),
						SKU:         "SKU-00002",
						ProductName: "Product Name 2",
						Quantity:    int64(3),
						UnitPrice:   float64(15000),
						TotalPrice:  float64(45000),
						Discount:    float64(2700),
						GrandTotal:  float64(42300),
					},
				},
				err: nil,
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mockRepo := foobarbaz_mock.NewMockFooRepository(ctrl)
				s := &foobarbaz.FooServiceImpl{
					FooRepository: mockRepo,
				}
				test.setupMock(mockRepo, test.entityID, *test.returns, *test.returnItems, test.err)
				got, err := s.ResolveByID(test.entityID, true)

				assert.Equal(t, test.err, err)
				assert.Equal(t, test.returns.Name, got.Name)
				assert.Equal(t, len(*test.returnItems), len(got.Items))
			})
		}
	})
}
