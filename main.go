package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var tpl *template.Template

func main() {
	// HTTP sunucusunu başlat
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	fmt.Println("HTTP server running at: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html"))
}
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method == "POST" {
			r.ParseForm()
			text := r.FormValue("ascii-data")
			font := r.FormValue("fonts")
			if hasTurkishChars(text) || hasTurkishChars(font) {
				http.Error(w, "404-Bad Request", http.StatusBadRequest)
				return
			}
			// ASCII sanatını oluştur
			asciiArt, err := generateASCIIArt(text, font)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// HTML sayfasına ASCII sanatını gönder
			tpl.ExecuteTemplate(w, "index.html", struct{ First string }{asciiArt})
			return
		}
		// İstek metodu GET ise ana sayfayı göster
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	} else if r.URL.Path == "/ascii-art" {
		if r.Method == "POST" {
			r.ParseForm()
			text := r.FormValue("ascii-data")
			font := r.FormValue("fonts")
			if hasTurkishChars(text) || hasTurkishChars(font) {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			// ASCII sanatını oluştur
			asciiArt, err := generateASCIIArt(text, font)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// HTML sayfasına ASCII sanatını gönder
			tpl.ExecuteTemplate(w, "index.html", struct{ First string }{asciiArt})
			return
		}
		// İstek metodu GET ise ASCII sanatı oluşturma sayfasını göster
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}
	// Diğer URL'lerde 404 Not Found hatası döndür
	http.NotFound(w, r)
}
func hasTurkishChars(s string) bool {
	turkishChars := "çÇğĞıİöÖşŞüÜ"
	for _, char := range s {
		if strings.ContainsRune(turkishChars, char) {
			return true
		}
	}
	return false
}
func generateASCIIArt(text, style string) (string, error) {
	// ASCII sanat dosyasını oku
	lines, err := readASCIIArt(style + ".txt")
	if err != nil {
		return "", err
	}
	// Text'teki başındaki ve sonundaki boşlukları kaldır
	text = strings.TrimSpace(text)
	var words []string
	if strings.Contains(text, "\\n") {
		words = strings.Split(text, "\\n")
	} else {
		words = strings.Split(text, "\n")
	}
	var result strings.Builder
	// Her bir kelime için ASCII sanatını oluştur
	for _, word := range words {
		if word == "" {
			result.WriteString("\n") // Boş satır ekle
			continue
		}
		result.WriteString(createWordASCII(word, lines)) // Kelimeyi ASCII sanatına dönüştür ve sonucu birleştir
		result.WriteString("\n")                         // Her kelimenin sonuna yeni satır ekle
	}
	return result.String(), nil
}

// ASCII sanat dosyasını satır bazında oku
func readASCIIArt(filename string) ([]string, error) {
	// Dosyayı aç
	file, err := os.Open("fonts/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	// Satırları tarayıcı kullanarak oku
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	// Hata kontrolü
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Her kelime için ASCII sanatını oluştur
func createWordASCII(word string, lines []string) string {
	var result strings.Builder
	// Her bir satır için
	for i := 1; i <= 8; i++ {
		// Kelimenin her bir karakteri için
		for _, char := range word {
			// Karakterin ASCII değerini al ve uygun satırı bul
			charIndex := int(char) - 32
			if charIndex < 0 || charIndex >= len(lines) {
				continue
			}
			lineIndex := (charIndex * 9) + i
			if lineIndex < 0 || lineIndex >= len(lines) {
				continue
			}
			// Bulunan satırı yazdır
			result.WriteString(lines[lineIndex])
		}
		result.WriteString("\n") // Kelimenin sonuna yeni satır ekle
	}
	return result.String()
}
