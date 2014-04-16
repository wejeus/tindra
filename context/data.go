package context

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"os"
	"path/filepath"
)

type Data map[string]interface{}

func NewData(path string) *Data {
	// data := Data{} // works
	var data Data = make(Data) // works
	// var data Data // do not work (nil map when using)
	// var data Data = new(Data) // do not work (pointer assignment)

	fi, _ := os.Stat(path)
	if !fi.IsDir() {
		return &data
	}

	dataFiles, _ := readFiles(path, map[string]bool{"yaml": true, "yml": true})

	// TODO: Is interface{} a pointer or will this cause an unnecessary copy?
	for filename, yamlData := range dataFiles {
		fmt.Println("parsing datafile: " + filename)
		var values interface{}
		yaml.Unmarshal(yamlData, &values)

		extension := filepath.Ext(filename)
		name := filename[0 : len(filename)-len(extension)]
		data[name] = values
		// fmt.Printf("%#v\n", values)
	}

	// docs := data["docs"]
	// fmt.Printf("%#v\n", docs)
	// fmt.Printf("%v\n", docs["title"])
	// for k, _ := range docs.([]interface{}) {
	// 	fmt.Printf("%#v", k)
	// }

	// showDATA(docs)

	return &data
}

func showDATA(d interface{}) {

	// var indentLevel = 0
	// indent := func() string {
	// 	var level string
	// 	for i := 0; i < indentLevel; i++ {
	// 		level += "\t"
	// 	}
	// 	return level
	// }

	var walker func(data interface{})

	walker = func(data interface{}) {

		switch outer := data.(type) {
		case string:
			fmt.Printf("string\n")
		case []interface{}:
			fmt.Printf("([]interface{})")
			for _, typed := range outer {
				walker(typed)
			}
		case map[string]interface{}:
			fmt.Printf("map[string]interface{}\n")

		case map[interface{}]interface{}:
			for innerKey, innerValue := range outer {
				switch inner := innerKey.(type) {
				case string:
					fmt.Printf("(inner) string\n")
				default:
					fmt.Printf("(inner) unknown: %T %T\n", inner, innerValue)
				}
			}
		default:
			fmt.Printf("unknown: %T\n", outer)
		}

		// if v, ok := data.(string); ok {
		// 	single := string(v)
		// 	fmt.Printf(indent()+"%s\n", single)
		// 	return
		// }

		// if v, ok := data.([]interface{}); ok {
		// 	for _, parentValue := range v {
		// 		walker(parentValue)
		// 	}
		// 	return
		// }

		// if v, ok := data.(map[interface{}]interface{}); ok {
		// 	for interfaceKey, interfaceValue := range v {
		// 		v1, isKeyString := interfaceKey.(string)
		// 		v2, isValueString := interfaceValue.(string)
		// 		v3, isValueStringArray := interfaceValue.([]string)
		// 		if isKeyString && isValueString {
		// 			fmt.Printf(indent()+"%s : %s\n", v1, v2)
		// 		} else if isKeyString && isValueStringArray {
		// 			fmt.Print(indent() + v1 + " : [")
		// 			for _, v := range v3 {
		// 				fmt.Print(v + " ")
		// 			}
		// 			fmt.Println("")
		// 		} else if isKeyString {
		// 			fmt.Println(indent() + v1 + " : [")
		// 			indentLevel++
		// 			walker(interfaceValue)
		// 			indentLevel--
		// 			fmt.Println(indent() + "]")
		// 		} else {
		// 			fmt.Printf("%T\n", interfaceKey)
		// 		}
		// 	}
		// 	return
		// }
	}

	walker(d)
}

// func showDATA(d interface{}) {

// 	var indentLevel = 0
// 	indent := func() string {
// 		var level string
// 		for i := 0; i < indentLevel; i++ {
// 			level += "\t"
// 		}
// 		return level
// 	}

// 	var walker func(data interface{})

// 	walker = func(data interface{}) {

// 		if v, ok := data.(string); ok {
// 			single := string(v)
// 			fmt.Printf(indent()+"%s\n", single)
// 			return
// 		}

// 		if v, ok := data.([]interface{}); ok {
// 			for _, parentValue := range v {
// 				walker(parentValue)
// 			}
// 			return
// 		}

// 		if v, ok := data.(map[interface{}]interface{}); ok {
// 			for interfaceKey, interfaceValue := range v {
// 				v1, isKeyString := interfaceKey.(string)
// 				v2, isValueString := interfaceValue.(string)
// 				v3, isValueStringArray := interfaceValue.([]string)
// 				if isKeyString && isValueString {
// 					fmt.Printf(indent()+"%s : %s\n", v1, v2)
// 				} else if isKeyString && isValueStringArray {
// 					fmt.Print(indent() + v1 + " : [")
// 					for _, v := range v3 {
// 						fmt.Print(v + " ")
// 					}
// 					fmt.Println("")
// 				} else if isKeyString {
// 					fmt.Println(indent() + v1 + " : [")
// 					indentLevel++
// 					walker(interfaceValue)
// 					indentLevel--
// 					fmt.Println(indent() + "]")
// 				} else {
// 					fmt.Printf("%T\n", interfaceKey)
// 				}
// 			}
// 			return
// 		}
// 	}

// 	walker(d)
// }
