package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func GenerateToken() (string, error) {
	const tokenSize = 32 // 32 bytes = 64 caracteres hex

	bytes := make([]byte, tokenSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error generating validation token: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func Base64ToString(entrada string) string {
	texto, err := base64.StdEncoding.DecodeString(entrada)
	if err != nil {
		fmt.Errorf("Error b64: ", err)
		return entrada
	}
	return string(texto)
}

func StringToB64(texto string) string {
	origen := []byte(fmt.Sprintf("%s", texto))
	destino := make([]byte, base64.StdEncoding.EncodedLen(len(texto)))
	base64.StdEncoding.Encode(destino, origen)
	return string(destino)
}

func ISOtoUTF(texto string) string {
	transformer := charmap.ISO8859_1.NewDecoder()
	utf, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(texto), transformer))
	if err != nil {
		return texto
	}
	//fmt.Printf("ISOtoUTF: [%s] en [%s]\n", texto, utf)
	return string(utf)
}

func RecibeDatosPost(r *http.Request, camposB64 []string) (map[string]string, int) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer r.Body.Close()

	var inputData map[string]interface{}
	err = json.Unmarshal(body, &inputData)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	data := make(map[string]string)
	for key, value := range inputData {
		esBase64 := false
		for _, campo := range camposB64 {
			if campo == strings.ToLower(key) {
				esBase64 = true
				break
			}
		}
		strValue, ok := value.(string)
		if !ok {
			strValue = fmt.Sprintf("%v", value)
		}
		if esBase64 {
			decodedValue, err := base64.StdEncoding.DecodeString(strValue)
			if err != nil {
				return nil, http.StatusBadRequest
			}
			data[strings.ToLower(key)] = ISOtoUTF(string(decodedValue))
		} else {
			data[strings.ToLower(key)] = strValue
		}
	}

	return data, 0
}

func RecibeDatosPut(r *http.Request, camposB64 []string) (map[string]string, int) {
	if r.Method != http.MethodPut {
		return nil, http.StatusMethodNotAllowed
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer r.Body.Close()

	var inputData map[string]interface{}
	err = json.Unmarshal(body, &inputData)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	data := make(map[string]string)
	for key, value := range inputData {
		esBase64 := false
		for _, campo := range camposB64 {
			if campo == key {
				esBase64 = true
				break
			}
		}
		strValue, ok := value.(string)
		if !ok {
			strValue = fmt.Sprintf("%v", value)
		}
		if esBase64 {
			decodedValue, err := base64.StdEncoding.DecodeString(strValue)
			if err != nil {
				return nil, http.StatusBadRequest
			}
			data[strings.ToLower(key)] = ISOtoUTF(string(decodedValue))
		} else {
			data[strings.ToLower(key)] = strValue
		}
	}

	return data, 0
}

func DevolverError(w http.ResponseWriter, error int) {
	if error == http.StatusMethodNotAllowed {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if error == http.StatusForbidden {
		http.Error(w, "Acción no permitida", http.StatusForbidden)
		return
	}

	if error == http.StatusInternalServerError {
		http.Error(w, "Error al procesar los datos", http.StatusInternalServerError)
		return
	}
	if error == http.StatusForbidden {
		http.Error(w, "El usuario no tiene los privilegios necesarios", http.StatusForbidden)
		return
	}

	http.Error(w, "Error interno", http.StatusInternalServerError)
}

func ToB64(texto string) string {
	origen := []byte(fmt.Sprintf("%s", texto))
	destino := make([]byte, base64.StdEncoding.EncodedLen(len(texto)))
	base64.StdEncoding.Encode(destino, origen)
	return string(destino)
}
