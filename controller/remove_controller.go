package controller

import (
	"errors"
	"fmt"
	"upload/bucket"
)

func RemoveAllBucketFiles() error {

	operation, err := bucket.LoadConfig()
	if err != nil {
		fmt.Println("Erro ao carregar configuração:", err)
		return err
	}
	removed, error := operation.RemoveBucketFiles()
	if error != nil {
		fmt.Println("Erro ao remover arquivos do bucket:", error)
		return err
	}
	if removed {
		fmt.Println("Arquivos removidos com sucesso.")
		return nil
	} else {
		fmt.Println("Nenhum arquivo removido.")
		return errors.New("Nenhum arquivo removido.")
	}

}
