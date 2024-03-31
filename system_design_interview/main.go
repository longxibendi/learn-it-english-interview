package main

import (
	"bufio"
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/kljensen/snowball"
)

var (
	keyboardControlChan = make(chan int, 100)
)

func main() {
	//speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.MPlayer{}}
	speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
	speech.Speak("To be a better guy,nice to meet you .")

	// get the dictionary
	dictWord := getDict("dict\\cmudict-0.7b-ipa.txt")
	// open the read book
	f, err := os.Open("system_design_interview\\system-design.txt")
	if err != nil {
		panic(err)
	}
	//color print
	colorPrint := color.New(color.FgCyan).Add(color.Bold)
	// control back 、go on 、pause
	// 'j' is back ,'k' is pause, 'l' is go on
	//keyboardControlChan := make(chan int ,100)

	go controlKeyboard(keyboardControlChan)
	//var lines []string
	sentenceSlice := make([]string, 0)
	r := bufio.NewReader(f)
	for {
		bytes, err := r.ReadBytes(byte('.'))
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
			//return lines, err
		}
		sentenceSlice = append(sentenceSlice, string(bytes))
	}

	for i := 0; i < len(sentenceSlice); i++ {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.

		if len(keyboardControlChan) > 0 {
			flowTag := <-keyboardControlChan
			if flowTag == 1 {
				//back to 3 sentences
				if i > 3 {
					i = i - 3
					continue
				}
			}
			if flowTag == 3 {
				// skip over 3 sentences
				if i < 1000000-3 {
					i = i + 3
					continue
				}
			}
			//if
		}
		////r.ReadLine()
		//bytes, err := r.ReadBytes(byte('.'))
		////bytes, _, err := r.ReadLine()
		//if err == io.EOF {
		//	break
		//}
		//if err != nil {
		//	panic(err)
		//	//return lines, err
		//}

		fmt.Printf("\n--------------------%d------------------------------\n", i)
		bytes := sentenceSlice[i]
		//fmt.Printf("%s\n", string(bytes))
		colorPrint.Printf("%s\n", string(bytes))
		phonetic2(string(bytes), dictWord)

		speech.Speak(string(bytes))
		time.Sleep(100 * time.Millisecond)

	}
}

func getDict(filePath string) map[string][]string {
	dictMap := make(map[string][]string)
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		//bytes, err := r.ReadBytes(byte('.'))
		bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
			//return lines, err
		}
		line := string(bytes)
		words := strings.Split(line, "\t")

		firstWord := words[0]
		phonetic := []string{}
		for j := 1; j < len(words); j++ {
			phonetic = append(phonetic, fmt.Sprintf("[%s]", words[j]))
		}
		dictMap[firstWord] = phonetic

	}
	return dictMap
}

func phonetic2(englishSentence string, dictWord map[string][]string) {
	// 要处理的英文句子
	//englishSentence := "Hello, how are you?"

	englishSentence = strings.ToUpper(englishSentence)
	//fmt.Println("englishSentence=" + englishSentence)
	// 将句子拆分为单词
	englishSentence = strings.ReplaceAll(englishSentence, "\n", "")
	words := strings.Split(englishSentence, " ")
	fmt.Printf("---------------------------------\n")
	// 遍历每个单词并获取音标
	for _, word := range words {
		// 使用 whatlanggo 包进行语言检测
		//info := whatlanggo.Detect(word)

		// 如果检测到的语言是英语，则提取音标
		//if info.Lang == whatlanggo.Eng {
		//fmt.Printf(" [%s]", word)
		phonetics, ok := dictWord[word]
		if !ok {
			continue
		}
		fmt.Printf("%s", phonetics[0])
		//} else {
		//fmt.Printf("%s", word)
		//fmt.Printf("*****************")
		//continue
		//log.Printf("无法获取单词 %s 的音标，因为它不是英语单词\n", word)
		//}
	}
}

func phonetic(englishSentence string) {
	// 要处理的英文句子
	//englishSentence := "Hello, how are you?"

	// 将句子拆分为单词
	words := strings.Split(englishSentence, " ")

	// 遍历每个单词并获取音标
	for _, word := range words {
		// 使用 snowball 包进行音标转换
		phonetic, err := snowball.Stem(word, "english", true)
		if err != nil {
			log.Fatalf("无法获取单词 %s 的音标：%v\n", word, err)
		}

		fmt.Printf("%s [%s]\n", word, phonetic)
	}
}

//func LeftClick() {
//	// 捕获鼠标点击事件
//	events := robotgo.Start()
//	defer robotgo.End()
//
//	// 标记左键是否按下
//	leftButtonPressed := false
//
//	for {
//		// 从事件通道中读取事件
//		event := <-events
//
//		// 检查事件类型是否为鼠标左键点击
//		if event.Kind == robotgo.MouseLeft {
//			if event.Type == robotgo.MouseDown {
//				leftButtonPressed = true
//			} else if event.Type == robotgo.MouseUp {
//				leftButtonPressed = false
//			}
//		}
//
//		// 如果左键一直按下，则输出"success"
//		if leftButtonPressed {
//			fmt.Println("success")
//		}
//	}
//}

func controlKeyboard(keyboardControlChan chan int) {
	// 打开键盘监听
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// 捕获连续按下的"j"字符
	count := 0
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		if char == 'j' {
			count++
			if count >= 10 {
				keyboardControlChan <- 1
				fmt.Println("连续按下了'j'字符")
				count = 0
			}
		} else if char == 'l' {
			count++
			if count >= 10 {
				keyboardControlChan <- 3
				fmt.Println("连续按下了'l'字符")
				count = 0
			}
		} else {
			count = 0
		}
	}
}
