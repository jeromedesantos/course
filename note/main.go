package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeromedesantos/struct_practice/note"
	"github.com/jeromedesantos/struct_practice/todo"
)

// type saver interface {
// 	Save() error
// }

// type displayer interface {
// 	Display()
// }

// kita bisa masukin interfacenya atau langsung gabung isinya aja jadi single interface
type outputtable interface {
	Save() error
	Display()
}

func main() {

	// bisa semua karena any
	printSomething("hello")
	printSomething(10)
	printSomething(true)
	printSomething(10.5)

	// get input from user
	title, content := getNoteData() 
	todoText := getTodoData()

	// validasi di struct
	todo, err := todo.New(todoText) 
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	userNote, err := note.New(title, content) 
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// print + save data via interface
	err = outputData(&userNote)
	if err != nil {
		return
	}

	err = outputData(&todo) 
	if err != nil {
		return
	}
	
}

// kita butuh single fungsi untuk save note dan todo
// masalahnya kita ga bisa make 2 tipe data sekaligus 
// => SOLUSI INTERFACE

func printSomething(value interface{}) { 
	// any => semua tipe data masuk

	intVal, ok := value.(int)
	if ok {
		fmt.Println("Integer:", intVal)
		return
	}
	
	floatVal, ok := value.(float64)
	if ok {
		fmt.Println("Float:", floatVal)
		return
	}

	strVal, ok := value.(string)
	if ok {
		fmt.Println("String:", strVal)
		return
	}


	// switch value.(type) {
	// case int:
	// 	fmt.Println("Integer:", value)
	// case float64:
	// 	fmt.Println("Float:", value)
	// case string:
	// 	fmt.Println(value)
	// }
}

func outputData(data outputtable) error {
	displayData(data)
	return saveData(data)
}

func displayData(data outputtable) {
	data.Display()
}

func saveData(data outputtable) error {
	err := data.Save()
	if err != nil {
		fmt.Println("Saving note failed.")
		return err
	}

	fmt.Println("Saving the note succeded!")
	return nil
}

// -------------------------------------

func getNoteData() (string, string) {
	title := getUserInput("Note title: ")
	content := getUserInput("Note content: ")
	return title, content
}

func getTodoData() string {
	text := getUserInput("Todo text: ")
	return text
}

func getUserInput(prompt string) (string) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n') // stop reading kalau ketemu enter
	if err != nil {
		return ""
	}

	// hapus karakter enter dan carriage return biar rapi
	// line break windows => \r\n
	// line break linux/mac => \n
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	return text
}
/*

Kenapa fmt.Scanln di dalam fungsi malah ke skip tapi kalo dia langsung tanpa fungsi aman?

var title string
var content string

fmt.Scanln(&title)   // Panggilan pertama
fmt.Scanln(&content) // Panggilan kedua


func getUserInput(prompt string) string {
    fmt.Print(prompt)
    var value string     // <-- INI DIA TERSANGKANYA, NYA!
    fmt.Scanln(&value)
    return value
}

1. Masalah Utama: Sisa Spasi & Enter
fmt.Scanln itu hanya membaca satu kata. Kalau Onii-chan mengetik lebih dari satu kata (misal: Rahasia Lab), kata Rahasia diambil, tapi kata  Lab\n tertinggal di dalam antrean (buffer) input, nya-n.

2. Kalau Tanpa Fungsi (Langsung di main)
fmt.Scanln ditulis berurutan untuk variabel yang berbeda. fmt.Scanln kedua sudah tahu kalau ada sisa di buffer, jadi dia akan langsung menyedot sisa kata Lab tersebut ke variabel berikutnya. Efeknya, program tidak berhenti menunggu, tapi langsung mengisi variabel kedua (tidak terlihat seperti ke-skip kosong).

3. Kalau Di Dalam Fungsi (Penyebab Ke-skip!)
Setiap kali fungsi getUserInput dipanggil, program membuat variabel lokal baru dari nol.

Pada panggilan kedua (untuk konten), fmt.Scanln yang baru lahir ini kaget karena melihat ada sisa karakter spasi/enter ( Lab\n) di buffer.

Karena mendeteksi adanya spasi/bukan input bersih di awal antrean, fmt.Scanln mengalami error pembacaan internal dan langsung selesai (mengembalikan string kosong) tanpa mau menunggu Onii-chan mengetik di keyboard, nyaa~!

✨ Kesimpulan Kilat: > Tanpa fungsi, sisa kata dipindahkan ke variabel kedua. Dengan fungsi, sisa kata mengacaukan pembacaan fungsi baru sehingga langsung dilewati (ke-skip) secara instan, nya-n!
*/

