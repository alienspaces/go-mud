package salesforce

import "encoding/xml"

// ResponseEnvelope -
type ResponseEnvelope struct {
	XMLName xml.Name     `xml:"soapenv:Envelope"`
	NsSoap  string       `xml:"xmlns:soapenv,attr"`
	NsXsd   string       `xml:"xmlns:xsd,attr"`
	NsXsi   string       `xml:"xmlns:xsi,attr"`
	Body    ResponseBody `xml:"soapenv:Body"`
}

// ResponseBody -
type ResponseBody struct {
	MessagesResponse MessagesResponse `xml:"messagesResponse"`
}

// MessagesResponse -
type MessagesResponse struct {
	XMLName xml.Name `xml:"http://soap.sforce.com/2005/09/outbound messagesResponse"`
	Ack     bool     `xml:"Ack"`
}

// ResponseEnvelopeUm - response envelope replacement for unmarshal in tests
type ResponseEnvelopeUm struct {
	XMLName xml.Name     `xml:"Envelope"`
	Body    ResponseBody `xml:"Body"`
}

// Response returns an ack or nack payload. Service Cloud will retry message delivery if the response is nack.
func Response(ack bool) *ResponseEnvelope {
	return &ResponseEnvelope{
		NsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		NsXsd:  "http://www.w3.org/2001/XMLSchema",
		NsXsi:  "http://www.w3.org/2001/XMLSchema-instance",
		Body: ResponseBody{
			MessagesResponse: MessagesResponse{
				Ack: ack,
			},
		},
	}
}
