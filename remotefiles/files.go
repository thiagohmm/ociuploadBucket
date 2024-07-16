package remotefiles

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Files struct {
	Src string
	Dst string
}

// copyFile copia um único arquivo de src para dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Copia os atributos do arquivo (permissões, timestamps, etc.)
	err = os.Chmod(dst, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir copia todos os arquivos da pasta src para a pasta dst
func (f Files) CopyDir() error {
	// Certifica-se de que a pasta de destino existe
	err := os.MkdirAll(f.Dst, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar diretórios para '%s': %w", filepath.Dir(f.Dst), err)
	}

	err = filepath.Walk(f.Src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erro ao copiar arquivo '%s': %w", path, err)
		}

		// Ignora diretórios
		if info.IsDir() {
			return nil
		}

		// Calcula o caminho de destino
		relPath, err := filepath.Rel(f.Src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(f.Dst, relPath)

		// Cria o diretório pai se não existir
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return err
		}

		// Copia o arquivo
		err = copyFile(path, destPath)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
