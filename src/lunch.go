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
	Similar string
	LookupCount int
	Acceptations []Acceptation
}

var question_re=regexp.MustCompile(`(?m)<div id="question"(\w|\W)*?</div>`)

var ul_re=regexp.MustCompile(`<ul(\w|\W)*?</ul>`)

var group_pos_re =regexp.MustCompile(`(?m)<div class="group_pos">(\w|\W)*?</div>`)
var pos_re =regexp.MustCompile(`(?m)<p>(\w|\W)*?</p>`)
var fl_re = regexp.MustCompile(`(?m)<strong class="fl">.*?</strong>`)
var label_re=regexp.MustCompile(`<span class="label_list">(\w|\W)*?</span>`)

var html_re = regexp.MustCompile("\\<[\\S\\s]+?\\>")  
var space_re =regexp.MustCompile(`(\s|&nbsp;){2,}`)

var history map[string]Word = make(map[string]Word)

var colorEnd string="\x1b[0m"

func printWithColor(value string,colorCode string,newline bool){
	fmt.Printf(colorCode)
	fmt.Printf(value)
	fmt.Printf(colorEnd)
	if newline{
		fmt.Printf("\n")
	}
}

// BIRed='\e[1;91m'        # Red
func red(value string,newline bool){
	printWithColor(value,"\x1b[1;91m",newline)
}

// BIGreen='\e[1;92m'      # Green
func green(value string,newline bool){
	printWithColor(value,"\x1b[1;92m",newline)
}
// BIYellow='\e[1;93m'     # Yellow
func yellow(value string,newline bool){
	printWithColor(value,"\x1b[1;93m",newline)
}
// BIBlue='\e[1;94m'       # Blue
func blue(value string,newline bool){
	printWithColor(value,"\x1b[1;94m",newline)
}
// BIPurple='\e[1;95m'     # Purple
func purple(value string,newline bool){
	printWithColor(value,"\x1b[1;95m",newline)
}
// BICyan='\e[1;96m'       # Cyan
func cyan(value string,newline bool){
	printWithColor(value,"\x1b[1;96m",newline)
}

func showError(err error){
	red(err.Error(),true)
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
		contents_str :=string(contents)
		question:=question_re.FindAllString(contents_str,-1)
		if len(question)>0{
			ul:=ul_re.FindAllString(question[0],-1)
			if len(ul)>0{
				similar:=html_re.ReplaceAllString(ul[0],"")
				similar=space_re.ReplaceAllString(similar," ")
				result.Similar=strings.TrimSpace(similar)
			}
		}
		group_pos:=group_pos_re.FindAllString(contents_str,-1)
		if len(group_pos)>0{
			ps :=pos_re.FindAllString(group_pos[0],-1)
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
		yellow("",true)
		notexist:=len(word.Similar)>0
		if notexist{
			sims:=strings.Split(word.Similar," ")
			yellow(">>>>>>>>The world "+word.Value+" not exist, do you means ",false)
			for _,val := range sims{
				cyan("["+val+"], ",false)
			}

			yellow("\b\b?\n",true)
			yellow("\tThe means of word ",false)
			cyan(sims[0],false)
			yellow(" is:\n",true)
		}
		for _,val := range word.Acceptations{
			if notexist{
				green("\t",false)				
			}
			green(val.PartOfSpeech,false)
			purple(val.Meaning,true)
		}
		yellow("",true)
	}
}

func main(){
	var word string
	bio:=bufio.NewReader(os.Stdin)
	for ;true;{
		word=""
		fmt.Printf("input your word/>")
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