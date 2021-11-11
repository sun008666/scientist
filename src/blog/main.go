package main

import "fmt"

func main(){
	flag := true
	map1 := map[string]int{
		"Jan":31,"Feb":28,"Mar":31,"Apr":30,"May":31,"jun":30,"Jul":31,"Aug":31,"Sep":30,"Oct":31,"Dec":31,
	}
	map2 := map[string]int{
		"Jan":31,"Feb":28,"Mar":31,"Apr":30,"May":31,"jun":30,"Jul":31,"Aug":31,"Sep":30,"Oct":31,"Dec":31,
	}
	if len(map1) != len(map2){
		flag = false
	}
	for k1,v1 := range map1{
		v2,ok := map2[k1]
		if ok {
			if v1 != v2 {
				flag = false
			}
		}else {
			flag = false
		}
	}
	fmt.Println(flag)

}
