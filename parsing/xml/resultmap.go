/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package xml

import "encoding/xml"

type IdArg struct {
	Column string `xml:"column,attr"`
	GoType string `xml:"type,attr"`
}

type Constructor struct {
	IdArg IdArg  `xml:"idArg"`
	Arg   string `xml:"arg"`
}

type Result struct {
	Property string `xml:"property,attr"`
	Column   string `xml:"column,attr"`
}

type ResultMap struct {
	XMLName xml.Name
	//id
	Id string `xml:"id,attr"`
	//struct类型名称
	TypeName string `xml:"type,attr"`
	//constructor - 用于在实例化类时，注入结果到构造方法中
	Constructor Constructor `xml:"constructor"`
	//一个 ID 结果；标记出作为 ID 的结果可以帮助提高整体性能
	ResultId Result `xml:"id"`
	//注入到字段或 Struct 属性的普通结果
	Results []Result `xml:"result"`
	//TODO:
	//association: 一个复杂类型的关联；许多结果将包装成这种类型
	//collection: 一个复杂类型的集合
	//discriminator: 使用结果值来决定使用哪个 resultMap
}
