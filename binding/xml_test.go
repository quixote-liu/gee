package binding

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testXmlDelivery struct {
	Name  string
	ID    string
	Email string
}

var xmlDelivery = testXmlDelivery{
	Name:  "lcs",
	ID:    "2020",
	Email: "123",
}

func TestXMLBind(t *testing.T) {
	xmlBytes, err := xml.Marshal(xmlDelivery)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodGet, "/xml/binding", bytes.NewReader(xmlBytes))

	xmlBinding := xmlBinding{}
	payload := testXmlDelivery{}
	err = xmlBinding.Bind(req, &payload)
	assert.Nil(t, err)

	assert.Equal(t, xmlDelivery, payload)
}

func TestXMLBindBody(t *testing.T) {
	xmlBytes, err := xml.Marshal(xmlDelivery)
	assert.Nil(t, err)

	xmlBinding := xmlBinding{}
	payload := testXmlDelivery{}
	err = xmlBinding.BindBody(xmlBytes, &payload)
	assert.Nil(t, err)

	assert.Equal(t, xmlDelivery, payload)
}
