package main

import "testing"

func TestMatchJSON(t *testing.T) {
	jsonData := `{ "key1":"value1", "key2":"value2" }`
	matcher := Matcher{
		Type:    "JSON",
		Pattern: "key2",
	}

	expect := "value2"
	result := matcher.match(jsonData)

	if result != expect {
		t.Errorf("expected %s, got %s", expect, result)
	}
}

func TestMatchURL(t *testing.T) {
	url := `http://host:8080/this/that/theother?foo=bar`
	matcher := Matcher{
		Type:    "URL",
		Pattern: "that",
	}

	expect := "that"
	result := matcher.match(url)

	if result != expect {
		t.Errorf("expected %s, got %s", expect, result)
	}
}

func TestMatchQueryString(t *testing.T) {
	qryStr := `http://host:8080/this/that/theother?foo=bar&this=that&theother`
	matcher := Matcher{
		Type:    "QueryString",
		Pattern: "this",
	}

	expect := "that"
	result := matcher.match(qryStr)

	if result != expect {
		t.Errorf("expected %s, got %s", expect, result)
	}

	matcher = Matcher{
		Type:    "QueryString",
		Pattern: "theother",
	}

	expect = "theother"
	result = matcher.match(qryStr)

	if result != expect {
		t.Errorf("expected %s, got %s", expect, result)
	}
}

func TestMatchXML(t *testing.T) {
	data := `
	<Data>
	<Person>
		<FullName>Grace R. Emlin</FullName>
		<Company>Example Inc.</Company>
		<Email where="home">
			<Addr>gre@example.com</Addr>
		</Email>
		<Email where='work' another='attribute'>
			<Addr>gre@work.com</Addr>
		</Email>
		<Group>
			<Value>Friends</Value>
			<Value>Squash</Value>
		</Group>
		<City>Hanga Roa</City>
		<State>Easter Island</State>
	</Person>
	</Data>
`
	matcher := Matcher{
		Type:    "XML",
		Pattern: "Data.Person.Email{where=work}.*Addr",
		// KeyPattern: "Data.Person.City",
	}

	expect := "gre@work.com"
	// expect := "Hanga Roa"
	result := matcher.match(data)

	if result != expect {
		t.Errorf("expected %s, got %s", expect, result)
	}
}
