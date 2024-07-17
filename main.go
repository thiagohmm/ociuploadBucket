package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"upload/controller"
	"upload/remotefiles"

	"github.com/joho/godotenv"
)

func init() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env:", err)
		return
	}

	pasta := "images"

	// Verifica se a pasta existe
	_, err = os.Stat(pasta)
	if os.IsNotExist(err) {
		// A pasta não existe, então cria
		err := os.MkdirAll(pasta, os.ModePerm)
		if err != nil {
			// Retorna um erro se não conseguir criar a pasta
			fmt.Errorf("erro ao criar a pasta %s: %v", pasta, err)
		}
		fmt.Printf("Pasta %s criada com sucesso.\n", pasta)
	} else if err != nil {
		// Retorna um erro se ocorrer um erro ao verificar a pasta
		fmt.Errorf("erro ao verificar a pasta %s: %v", pasta, err)
	} else {
		fmt.Printf("Pasta %s já existe.\n", pasta)
	}

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Escolha uma opção:")
		fmt.Println("1. Copy Images")
		fmt.Println("2. Delete Images Bucket")
		fmt.Println("3. Upload Bucket")
		fmt.Println("4. Sair")

		// Lê a opção do usuário
		scanner.Scan()
		opcao := scanner.Text()

		switch opcao {
		case "1":
			fmt.Println("Copiando imagens...")

			// Obtenha os caminhos de origem e destino do ambiente
			srcPath := os.Getenv("SRC_PATH")
			dstPath := "./images"

			fmt.Println("SRC_PATH:", srcPath)
			fmt.Println("DST_PATH:", dstPath)
			if srcPath == "" || dstPath == "" {
				fmt.Println("Por favor, defina SRC_PATH e DST_PATH no ambiente.")
				continue
			}

			fileToUpload := remotefiles.Files{Src: srcPath, Dst: dstPath}

			var wg sync.WaitGroup

			wg.Add(1) // Adiciona uma contagem para a WaitGroup
			go func() {
				defer wg.Done() // Decrementa a contagem na WaitGroup quando a goroutine terminar
				err := fileToUpload.CopyDir()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Imagens copiadas com sucesso!")
				}
			}()

			wg.Wait() // Aguarda a conclusão da goroutine
		case "2":
			fmt.Println("Deletando imagens do bucket...")

			err := controller.RemoveAllBucketFiles()
			if err != nil {
				fmt.Println(err)
			}

		case "3":
			fmt.Println("Fazendo upload para o bucket...")
			basePath := "./images"
			err := controller.UploadImages(basePath)
			if err != nil {
				fmt.Println(err)
			}

		case "4":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida. Por favor, escolha uma opção entre 1 e 4.")
		}
	}
}
