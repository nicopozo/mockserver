package main

import (
	"fmt"
	"regexp"
)

func main() {

	/*	s := ""

		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}

		fmt.Println(u.Scheme)
		fmt.Println(u.Host)
		fmt.Println(u.Path)
		fmt.Println(u.RawQuery)*/

	//path := "/single_player/{prepaid}/{id}/"
	//parts := strings.Split(path, ":")
	//strings.Replace()

	//var re = regexp.MustCompile(`/test/.+?/[1-9a-zA-z]+?$`)
	var re = regexp.MustCompile(`^/test/[^/]+?/[^/]+?/[^/]+$`)
	res := re.MatchString("/test/1/2")
	fmt.Println(res)

}

/*
func cmp_any(obj1, obj2 interface{}, op string) (bool, error) {
	switch op {
	case "<", "<=", "==", ">=", ">":
	default:
		return false, fmt.Errorf("op should only be <, <=, ==, >= and >")
	}
	fmt.Println("cmp_any: ", obj1, obj2)
	exp := fmt.Sprintf("%v %s %v", obj1, op, obj2)
	fmt.Println("exp: ", exp)
	fset := token.NewFileSet()
	res, err := types.Eval(fset, nil, 0, exp)
	if err != nil {
		return false, err
	}
	if res.IsValue() == false || (res.Value.String() != "false" && res.Value.String() != "true") {
		return false, fmt.Errorf("result should only be true or false")
	}
	if res.Value.String() == "true" {
		return true, nil
	}
	return false, nil
}
*/
