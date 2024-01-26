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

	"github.com/kljensen/snowball"
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
	//var lines []string
	r := bufio.NewReader(f)
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		bytes, err := r.ReadBytes(byte('.'))
		//bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
			//return lines, err
		}

		fmt.Printf("\n--------------------------------------------------\n")
		fmt.Printf("%s\n", string(bytes))
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
