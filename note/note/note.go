package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (note *Note) Display () {
	fmt.Printf("Your note titled %v has the following content:\n\n%v\n", note.Title, note.Content)
}

// tidak perlu pointer karena cuma simpan copyan data
func (note Note) Save() error {

	fileName := note.Title
	
	fileName = strings.ReplaceAll(fileName, " ", "_") // spasi jadi underscore
	fileName = strings.ToLower(fileName) + ".json" // convert semua uppercase jadi lowercase


	json, err := json.Marshal(note) // konversi ke json

	if err != nil {
		return err
	}

	return os.WriteFile(fileName , json, 0644)
	// 0644 => 3 digit pertama permission untuk owner, 3 digit berikutnya untuk group, 3 digit terakhir untuk all
	// 6 => rw (read write)
	// 4 => r (read)
	// jadi artinya owner bisa baca tulis, sisanya cuma bisa baca doang
}

// akses struct note pake ini
func New(title, content string) (Note, error) {
	if title == "" || content == "" {
		return Note{}, errors.New("Invalid input.")
	}

	return Note{
		Title:   title,
		Content: content,
		CreatedAt: time.Now(),
	}, nil
}