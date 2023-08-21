package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	// Константа регулярного выражения, по которой выполняется поиск ссылок
	uriRegex = `(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`
)

func main() {
	// Чтение аргументов командной строки
	args := cmdArgs{}
	err := args.parseArgs()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	// Скачивание вебсайта
	err = getWebsite(args)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

// Функция для скачивания сайта
func getWebsite(args cmdArgs) error {
	// URI сайта
	uri := args.uri

	// Максимальная глубина рекурсии при скачивании сайта
	mapDepth := 0
	if args.recursive {
		mapDepth = args.depth
	}

	// Парсинг URI и получение названия сайта
	url, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}
	folder := "data/" + url.Host + "/"

	// Создание для сайта своей папки
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	// рекурсивная загрузка сайта
	if err := download(uri, folder, mapDepth); err != nil {
		return err
	}

	return nil
}

// Функция для скачаивания сайта по URI uri в папку folder размером рекурсии depth
func download(uri string, folder string, depth int) error {
	// Проверка uri
	if _, err := url.ParseRequestURI(uri); err != nil {
		return err
	}

	// Скачивание ресурса
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Чтение ресурса
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Получение названия ресурса
	filename, err := parseFilename(uri, res.Header)
	if err != nil {
		return err
	}

	// Создание файла ресурса
	out, err := os.Create(folder + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Запись данных ресурса
	out.Write(data)

	// Если достугнута максимальная глубина рекурсии, то выход
	if depth < 1 {
		return nil
	}

	// Получение ссылок со страницы на другие ресурсы
	links := getSublinks(data)

	// Итерация по ссылкам и загрузка ресурсов
	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)

		go func(link string) {
			defer wg.Done()
			download(link, folder, depth-1)
		}(link)
	}
	wg.Wait()

	return nil
}

// Структура флагов
type cmdArgs struct {
	// Включить рекурсивную загрузку
	recursive bool
	// Максимальная длина рекурсивной загрузки
	depth int
	// URI сайта
	uri string
}

// Метод для инициализирования аргументов
func (a *cmdArgs) parseArgs() error {
	// Привязка аргументов к полям структуры
	flag.BoolVar(&a.recursive, "r", false, "Включить рекурсивную загрузку")
	flag.IntVar(&a.depth, "l", 0, "Максимальная длина рекурсивной загрузки")
	flag.Parse()

	if !a.recursive && a.depth != 0 {
		return errors.New(`Flag "l" can be defined only with "r" flag`)
	}

	a.uri = flag.Arg(0)
	if len(a.uri) == 0 {
		return errors.New("Required field uri is missing")
	}

	return nil
}

// Получение названия ресурса
func parseFilename(uri string, header http.Header) (string, error) {
	contentType := header.Get("Content-Type")
	mimeType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	mediaType := mimeType

	if strings.Contains(mimeType, "/") {
		mediaType = strings.Split(mimeType, "/")[1]
	}

	filename := getNameByUrl(uri, mediaType)

	return filename, nil
}

// Получение всех ссылок из data
func getSublinks(data []byte) []string {
	re := regexp.MustCompile(uriRegex)
	links := re.FindAll(data, -1)

	subUris := make([]string, len(links))
	for _, link := range links {
		subUris = append(subUris, string(link))
	}
	return subUris
}

// Получение названия по URL url
func getNameByUrl(url string, mediaType string) string {
	n := path.Base(url)

	if mediaType == "" {
		return ""
	}

	name := n
	if strings.Contains(n, "#") {
		name = strings.Split(n, "#")[0]
	}

	if strings.Contains(name, "?") {
		name = strings.Split(name, "?")[0]
	}

	if path.Ext(name) == "" && mediaType != "" {
		return name + "." + mediaType
	}

	return name
}
