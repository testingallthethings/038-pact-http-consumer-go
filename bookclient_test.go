package book_test

import (
	"errors"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/suite"
	book "github.com/testingallthethings/038-pact-http-consumer-go"
	"net/http"
	"testing"
)

type BookClientSuite struct {
	suite.Suite
}

var (
	pact *dsl.Pact
)

func TestBookClientSuite(t *testing.T) {
	suite.Run(t, new(BookClientSuite))
}

func (s *BookClientSuite) SetupSuite() {
	pact = &dsl.Pact{
		Consumer: "MarksBookClient",
		Provider: "BookApi",
		PactDir:  "./pacts",
	}
}

func (s *BookClientSuite) TearDownSuite() {
	pact.Teardown()
}

func (s *BookClientSuite) TestGetBookThatDoesNotExist() {
	pact.AddInteraction().
		Given("There is not a book with ISBN 123456789").
		UponReceiving("A GET request for book with ISBN 123456789").
		WithRequest(
			dsl.Request{
				Method: http.MethodGet,
				Path:   dsl.String("/book/123456789"),
				Headers: dsl.MapMatcher{
					"Accept": dsl.String("application/json"),
				},
			},
		).
		WillRespondWith(
			dsl.Response{
				Status: http.StatusNotFound,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
				Body: dsl.Match(book.JsonError{}),
			},
		)

	test := func() error {
		c := book.NewClient(fmt.Sprintf("http://localhost:%d", pact.Server.Port))
		_, err := c.GetBook("123456789")

		if err.Error() != "1234 - No book with ISBN 123456789" {
			return errors.New("error returned not as expected")
		}

		return nil
	}

	s.NoError(pact.Verify(test))

}

func (s *BookClientSuite) TestGetBookThatDoesExist() {
	pact.AddInteraction().
		Given("Book with ISBN 987654321 exists").
		UponReceiving("A GET request for book with ISBN 987654321").
		WithRequest(
			dsl.Request{
				Method: http.MethodGet,
				Path:   dsl.String("/book/987654321"),
				Headers: dsl.MapMatcher{
					"Accept": dsl.String("application/json"),
				},
			},
		).
		WillRespondWith(
			dsl.Response{
				Status: http.StatusOK,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
				Body: dsl.Match(book.Book{}),
			},
		)

	test := func() error {
		c := book.NewClient(fmt.Sprintf("http://localhost:%d", pact.Server.Port))
		book, err := c.GetBook("987654321")

		s.NoError(err)
		s.Equal("987654321", book.ISBN)
		s.Equal("Testing All The Things", book.Title)
		s.Equal("testing.jpg", book.Image)
		s.Equal("Computers", book.Genre)
		s.Equal(2021, book.YearPublished)

		return nil
	}

	s.NoError(pact.Verify(test))
}
