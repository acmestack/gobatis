/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package xml

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"

	"github.com/acmestack/gobatis/logging"
)

const (
	MapperStart = "mapper"
)

func ParseFile(path string) (*Mapper, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		logging.Warn("error: %v", err)
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logging.Warn("error: %v", err)
		return nil, err
	}

	return Parse(data)
}

func Parse(data []byte) (*Mapper, error) {
	v := Mapper{}
	err := xml.Unmarshal(data, &v)
	if err != nil {
		logging.Warn("error: %v", err)
		return nil, err
	}
	return &v, nil
}

func parseInner(reader io.Reader) {
	decoder := xml.NewDecoder(reader)
	var strName string
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		name := getStartElementName(token)
		if name != nil {
			if name.Local == MapperStart {
				switch t := token.(type) {
				case xml.StartElement:
					startElement := xml.StartElement(t)
					logging.Debug("start: ", startElement.Name.Local)
					strName = startElement.Name.Local
				case xml.EndElement:
					endElement := xml.EndElement(t)
					logging.Debug("end: ", endElement.Name.Local)
				case xml.CharData:
					data := xml.CharData(t)
					str := string(data)
					switch strName {
					case "City":
						logging.Debug("city:", str)
					case "first":
						logging.Debug("first: ", str)
					}
				}
				break
			}
		}
	}
}

func getStartElementName(token xml.Token) *xml.Name {
	switch t := token.(type) {
	case xml.StartElement:
		startElement := xml.StartElement(t)
		logging.Debug("start: ", startElement.Name.Local)
		return &startElement.Name
	}
	return nil
}

//   * If the struct has a field of type []byte or string with tag
//      ",innerxml", Unmarshal accumulates the raw XML nested inside the
//      element in that field. The rest of the rules still apply.
//
//   * If the struct has a field named XMLName of type Name,
//      Unmarshal records the element name in that field.
//
//   * If the XMLName field has an associated tag of the form
//      "name" or "namespace-URL name", the XML element must have
//      the given name (and, optionally, name space) or else Unmarshal
//      returns an error.
//
//   * If the XML element has an attribute whose name matches a
//      struct field name with an associated tag containing ",attr" or
//      the explicit name in a struct field tag of the form "name,attr",
//      Unmarshal records the attribute value in that field.
//
//   * If the XML element has an attribute not handled by the previous
//      rule and the struct has a field with an associated tag containing
//      ",any,attr", Unmarshal records the attribute value in the first
//      such field.
//
//   * If the XML element contains character data, that data is
//      accumulated in the first struct field that has tag ",chardata".
//      The struct field may have type []byte or string.
//      If there is no such field, the character data is discarded.
//
//   * If the XML element contains comments, they are accumulated in
//      the first struct field that has tag ",comment".  The struct
//      field may have type []byte or string. If there is no such
//      field, the comments are discarded.
//
//   * If the XML element contains a sub-element whose name matches
//      the prefix of a tag formatted as "a" or "a>b>c", unmarshal
//      will descend into the XML structure looking for elements with the
//      given names, and will map the innermost elements to that struct
//      field. A tag starting with ">" is equivalent to one starting
//      with the field name followed by ">".
//
//   * If the XML element contains a sub-element whose name matches
//      a struct field's XMLName tag and the struct field has no
//      explicit name tag as per the previous rule, unmarshal maps
//      the sub-element to that struct field.
//
//   * If the XML element contains a sub-element whose name matches a
//      field without any mode flags (",attr", ",chardata", etc), Unmarshal
//      maps the sub-element to that struct field.
//
//   * If the XML element contains a sub-element that hasn't matched any
//      of the above rules and the struct has a field with tag ",any",
//      unmarshal maps the sub-element to that struct field.
//
//   * An anonymous struct field is handled as if the fields of its
//      value were part of the outer struct.
//
//   * A struct field with tag "-" is never unmarshaled into.
//
