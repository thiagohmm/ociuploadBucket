package controller

import (
	"fmt"
	"upload/bucket"
)

func ListBucketFiles() error {

	operation, err := bucket.LoadConfig()
	if err != nil {
		fmt.Println("Erro ao carregar configuração:", err)
		return err
	}
	error := operation.ListarBucket()
	if error != nil {
		fmt.Println("Erro ao listar arquivos do bucket:", error)
		return err
	}
	return nil

}
