package main

//Author: Sergio Cabrera Cirilo
import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func getUrl(name string) string {
	privateKey := "3a192df8736e821281f470036439061634649e23"
	publicKey := "c3bd0e5b37d3c0133716463fb7fe5dc9"
	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	h := md5.New()
	io.WriteString(h, ts+privateKey+publicKey)
	hash := hex.EncodeToString(h.Sum(nil))
	rest := "http://gateway.marvel.com/v1/public/characters?ts=" + ts + "&apikey=" + publicKey + "&hash=" + hash
	if name != "" {
		name = url.QueryEscape(name)
		rest = "http://gateway.marvel.com/v1/public/characters?name=" + name + "&ts=" + ts + "&apikey=" + publicKey + "&hash=" + hash
	}
	return rest
}

type Characters struct {
	Code int `json:"code"`
	Data struct {
		Results []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"results"`
	} `json:"data"`
}

func getCharacters(name string) (Characters, error) {
	url := getUrl(name)
	resp, err := http.Get(url)
	if err != nil {
		return Characters{}, err
	}
	defer resp.Body.Close()

	characters := Characters{}
	json.NewDecoder(resp.Body).Decode(&characters)
	if characters.Code != 200 {
		return Characters{}, errors.New("Ocurrió un error al obtener datos")
	}
	return characters, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n*******************************************************************")
		fmt.Print("Menú\n1. Buscar por nombre\n2. Listar\n3. Salir\nTeclea el número de la opción: ")
		opt, _ := reader.ReadString('\n')
		opt = strings.TrimRight(opt, "\r\n")
		switch opt {
		case "1":
			fmt.Print("Introduce el nombre: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimRight(name, "\r\n")
			characters, err := getCharacters(name)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(characters.Data.Results) > 0 {
				character := characters.Data.Results[0]
				fmt.Printf("\nID: %d\nNombre: %s\nDescripción: %s\n", character.ID, character.Name, character.Description)
			} else {
				fmt.Println("\nEl personaje no fué encontrado")
			}
			break
		case "2":
			characters, err := getCharacters("")
			if err != nil {
				fmt.Println(err)
				return
			}
			for i, character := range characters.Data.Results {
				fmt.Printf("#%d\nID: %d\nNombre: %s\nDescripción: %s\n\n", i+1, character.ID, character.Name, character.Description)
			}
			break
		default:
			return
		}
		fmt.Print("\nPresiona ENTER para continuar")
		reader.ReadBytes('\n')
	}
}
