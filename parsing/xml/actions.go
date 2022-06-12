/*
 * Copyright (c) 2022, OpeningO
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

import "encoding/xml"

type Select struct {
	XMLName       xml.Name
	Id            string `xml:"id,attr"`
	ParameterType string `xml:"parameterType,attr"`
	ParameterMap  string `xml:"parameterMap,attr"`
	ResultType    string `xml:"resultType,attr"`
	ResultMap     string `xml:"resultMap,attr"`
	FlushCache    string `xml:"flushCache,attr"`
	UseCache      string `xml:"useCache,attr"`
	Timeout       string `xml:"timeout,attr"`
	FetchSize     string `xml:"fetchSize,attr"`
	StatementType string `xml:"statementType,attr"`
	ResultSetType string `xml:"resultSetType,attr"`

	//If       []If    `xml:"if"`
	//Include Include `xml:"include"`
	//Where   Where   `xml:"where"`
	//Data    string  `xml:",chardata"`
	Data string `xml:",innerxml"`
}

type Insert struct {
	XMLName          xml.Name
	Id               string `xml:"id,attr"`
	ParameterType    string `xml:"parameterType,attr"`
	FlushCache       string `xml:"flushCache,attr"`
	Timeout          string `xml:"timeout,attr"`
	StatementType    string `xml:"statementType,attr"`
	UseGeneratedKeys string `xml:"useGeneratedKeys,attr"`
	KeyProperty      string `xml:"keyProperty,attr"`
	KeyColumn        string `xml:"keyColumn,attr"`

	//If       []If    `xml:"if"`
	//Include Include `xml:"include"`
	//Where   Where   `xml:"where"`
	//Data    string  `xml:",chardata"`
	Data string `xml:",innerxml"`
}

type Update struct {
	XMLName       xml.Name
	Id            string `xml:"id,attr"`
	ParameterType string `xml:"parameterType,attr"`
	FlushCache    string `xml:"flushCache,attr"`
	Timeout       string `xml:"timeout,attr"`
	StatementType string `xml:"statementType,attr"`

	//If       []If    `xml:"if"`
	//Include Include `xml:"include"`
	//Set     Set     `xml:"set"`
	//Where   Where   `xml:"where"`
	//Data    string  `xml:",chardata"`
	Data string `xml:",innerxml"`
}

type Delete struct {
	XMLName       xml.Name
	Id            string `xml:"id,attr"`
	ParameterType string `xml:"parameterType,attr"`
	FlushCache    string `xml:"flushCache,attr"`
	Timeout       string `xml:"timeout,attr"`
	StatementType string `xml:"statementType,attr"`

	//If       []If    `xml:"if"`
	//Include Include `xml:"include"`
	//Where   Where   `xml:"where"`
	//Data    string  `xml:",chardata"`
	Data string `xml:",innerxml"`
}

func (a *Select) ParseDynamic() {

}
