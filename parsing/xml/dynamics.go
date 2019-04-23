/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package xml

type Foreach struct {
    Item       string `xml:"item,attr"`
    Collection string `xml:"collection,attr"`
    Separator  string `xml:"separator,attr"`
    Index      string `xml:"index,attr"`
    Open       string `xml:"open,attr"`
    Close      string `xml:"close,attr"`
}

type Sql struct {
    Id  string `xml:"id,attr"`
    Sql string `xml:",chardata"`
}

type Property struct {
    Name  string `xml:"name,attr"`
    Value string `xml:"value,attr"`
}

type Include struct {
    Refid      string     `xml:"refid,attr"`
    Properties []Property `xml:"property"`
}

type If struct {
    Test string `xml:"test,attr"`
    If   string `xml:",chardata"`
}

type Where struct {
    If []If `xml:"if"`
}

type Set struct {
    If []If `xml:"if"`
}
