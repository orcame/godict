package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"regexp")

var pos_re =regexp.MustCompile(`(?m)<div class="group_pos">(\w|\W)*?</div>`)
var p_re =regexp.MustCompile(`(?m)<p>(\w|\W)*?</p>`)
var html_re = regexp.MustCompile("\\<[\\S\\s]+?\\>")  
var spanc_re =regexp.MustCompile(`\s{2,}`)


func translate(word string) {
	resp, err := http.Get("http://www.iciba.com/"+word)
	if err !=nil{
		fmt.Println("%s",err)
	}else{
		defer resp.Body.Close()
		contents, err:=ioutil.ReadAll(resp.Body)
		if err != nil{
			fmt.Printf("%s",err)
			os.Exit(1)
		}
		pos:=pos_re.FindAllString(string(contents),-1)
		ps :=p_re.FindAllString(string(pos[0]),-1)
		
		for _, value:= range ps{
			v:=html_re.ReplaceAllString(value,"")
			v=spanc_re.ReplaceAllString(v," ")
			fmt.Println(v)
		}
	}
}

func main(){
	var word string
	for ;true;{
		word=""
		fmt.Printf("input your word/>")
		fmt.Scanf("%s",&word)
		if len(word)>0{			
			translate(word)
		}
	}
}