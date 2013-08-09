package main

import (
	"fmt"
	"net/http"
	"bufio"
	"os"
	"strings"
	"io/ioutil"
	"regexp")

type Acceptation struct{
	PartOfSpeech string
	Meaning string
}

type Word struct{
	Value string
	LookupCount int
	Acceptations []Acceptation
}

var group_pos_re =regexp.MustCompile(`(?m)<div class="group_pos">(\w|\W)*?</div>`)
var pos_re =regexp.MustCompile(`(?m)<p>(\w|\W)*?</p>`)
var fl_re = regexp.MustCompile(`(?m)<strong class="fl">.*?</strong>`)
var label_re=regexp.MustCompile(`<span class="label_list">(\w|\W)*?</span>`)

var html_re = regexp.MustCompile("\\<[\\S\\s]+?\\>")  
var space_re =regexp.MustCompile(`\s{2,}`)

var history map[string]Word = make(map[string]Word)

func showError(err error){
	fmt.Printf("\x1b[1;91m%s\x1b[0m\n",err)
}

func showPartOfSpeech(val string){
	fmt.Printf("\x1b[1;92m"+val+"\x1b[0m")
}

func showMeaning(val string) {
	fmt.Printf("\x1b[1;95m"+val+"\x1b[0m")	
}


func translate(word string) Word{
	if val,ok := history[word];ok{
		val.LookupCount++
		return val
	}
	var result Word
	result.Value=word

	resp, err := http.Get("http://www.iciba.com/"+word) 	
	if err !=nil{
		showError(err)
		return result
	}else{

		defer resp.Body.Close()
		contents, err:=ioutil.ReadAll(resp.Body)
		if err != nil{
			fmt.Printf("%s",err)
		}
		group_pos:=group_pos_re.FindAllString(string(contents),-1)
		if len(group_pos)>0{
			ps :=pos_re.FindAllString(string(group_pos[0]),-1)
			pos_len:=len(ps)

			accs:=make([]Acceptation,pos_len)

			for idx, value:= range ps{
				// fmt.Printf("%d",idx)
				fl:=fl_re.FindAllString(value,-1)
				if len(fl)>0{
					pos_val:=html_re.ReplaceAllString(fl[0],"")
					pos_val = strings.Replace(pos_val,"&lt;","<",-1)
					pos_val = strings.Replace(pos_val,"&gt;",">",-1)
					accs[idx].PartOfSpeech=pos_val
				}
				label:=label_re.FindAllString(value,-1)
				if len(label)>0{
					meaning:=html_re.ReplaceAllString(label[0],"")
					meaning=space_re.ReplaceAllString(meaning," ")
					accs[idx].Meaning=meaning
				}
				result.Acceptations=accs
			}
		}
	}
	result.LookupCount++
	history[word]=result
	return result
}

func showWord(word Word){
	if word.LookupCount>0{
		for _,val := range word.Acceptations{
			showPartOfSpeech(val.PartOfSpeech)
			showMeaning(val.Meaning)
			fmt.Println()
		}
	}
}

func main(){
	var word string
	bio:=bufio.NewReader(os.Stdin)
	for ;true;{
		word=""
		fmt.Printf("\ninput your word/>")
		line ,_,err:=bio.ReadLine()
		if err !=nil{
			showError(err)
		}
		word=string(line)
		word=strings.TrimSpace(word)
		if len(word)>0{			
			val:=translate(word)
			showWord(val)
		}
	}
}