package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"upload/bucket"
)

func UploadImages(basePath string) error {
	var wg sync.WaitGroup

	var images []string
	uploadControl := make(chan struct{}, 100) // Canal de controle para limitar uploads simultâneos

	// Percorre os diretórios a partir do basePath para encontrar imagens
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Verifica se o arquivo tem uma extensão de imagem
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if ext == ".jpeg" || ext == ".jpg" || ext == ".png" || ext == ".gif" {
				images = append(images, path)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("erro ao percorrer diretórios: %w", err)
	}

	operation, err := bucket.LoadConfig()
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Adiciona uma tarefa ao WaitGroup para cada imagem
	for _, image := range images {
		wg.Add(1)
		go func(image string) {
			defer wg.Done() // Garante que Done será chamado ao final da goroutine

			// Reserva um espaço no canal de controle de upload
			uploadControl <- struct{}{}

			fmt.Println("Uploading:", image)
			objectName := image
			// Supondo que operation.UploadObject é uma função válida que faz upload do objeto
			result, err := operation.UploadObject(objectName)
			if err != nil {
				fmt.Println("Erro ao realizar upload do objeto:", err)
			} else {
				fmt.Println(result)
			}

			// Libera o espaço no canal de controle de upload
			<-uploadControl
		}(image)
	}

	wg.Wait() // Aguarda a conclusão de todas as goroutines

	return nil
}
