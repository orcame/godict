package main

import (
	"net/http"
	"bufio"
	"os"
	"strings"
	"io/ioutil"
	"../../terminal/src"
	"regexp")

type Acceptation struct{
	PartOfSpeech string
	Meaning string
}

type Word struct{
	Word string
	Real string
	Similar string
	LookupCount int
	Acceptations []Acceptation
}
var(
	to=terminal.Stdout.Bold().Intensity()
	te=terminal.Stderr.Color('r').Bold().Intensity()
	question_re=regexp.MustCompile(`(?m)<div id="question"(\w|\W)*?</div>`)

	ul_re=regexp.MustCompile(`<ul(\w|\W)*?</ul>`)
	word_re=regexp.MustCompile(`<h1 id="word_name_h1">.*?</h1>`)

	group_pos_re =regexp.MustCompile(`(?m)<div class="group_pos">(\w|\W)*?</div>`)
	pos_re =regexp.MustCompile(`(?m)<p>(\w|\W)*?</p>`)
	fl_re = regexp.MustCompile(`(?m)<strong class="fl">.*?</strong>`)
	label_re=regexp.MustCompile(`<span class="label_list">(\w|\W)*?</span>`)

	html_re = regexp.MustCompile("\\<[\\S\\s]+?\\>")  
	space_re =regexp.MustCompile(`(\s|&nbsp;){2,}`)

	history map[string]Word = make(map[string]Word)
)


func translate(word string) Word{
	if val,ok := history[word];ok{
		val.LookupCount++
		return val
	}
	var result Word=Word{Word:word}

	resp, err := http.Get("http://www.iciba.com/"+word) 	
	if err !=nil{
		te.Print(err.Error()).Nl()
		return result
	}else{

		defer resp.Body.Close()
		contents, err:=ioutil.ReadAll(resp.Body)
		if err != nil{
			te.Print(err.Error()).Nl()
			return result;
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
		word_r :=word_re.FindAllString(contents_str,-1)
		if len(word_r)>0{
			result.Real=html_re.ReplaceAllString(word_r[0],"")
		}else{
			to.Bold().Fprint("@{r}the word ","@{c}",word,"@{r} not exist.").Nl()
			return result;
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
		terminal.Stdout.Color('y').Nl()
		notexist:=len(word.Similar)>0
		if notexist{
			sims:=strings.Split(word.Similar," ")
			to.Fprint("@{y}>>>>>>>>The world ","@{c}",word.Word,"@{y} not exist, do you means ").
				Color('c')
			for _,val := range sims{
				to.Print("[",val,"], ")
			}
			to.Fprint("@{y}\b\b?").Nl().Fprint("@{y}\tThe means of word ","@{c}",word.Real,"@{y} is:").Nl()
		}
		for _,val := range word.Acceptations{
			if notexist{
				to.Color('g').Print("\t")				
			}
			to.Color('g').Print(val.PartOfSpeech).
				Color('c').Print(val.Meaning).Nl()
		}
		to.Nl()
	}
}

func main(){
	var word string
	if len(os.Args)>1{
		word = strings.Join(os.Args[1:]," ")
		val:=translate(word)
		showWord(val)
		to.Reset()
	}else{
		bio:=bufio.NewReader(os.Stdin)
		for{
			to.Color('w').Print("input your word/>")
			line ,_,err:=bio.ReadLine()
			if err !=nil{
				te.Print(err.Error()).Nl()
			}
			word=string(line)
			if(word == ":q"){
				to.Reset()
				break
			}
			word=strings.TrimSpace(word)
			if len(word)>0{			
				val:=translate(word)
				showWord(val)
			}
			word=""
		}
	}
}
