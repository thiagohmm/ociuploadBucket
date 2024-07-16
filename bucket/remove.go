package bucket

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

func (b BucketConfig) RemoveBucketFiles() (bool, error) {
	keyFileContent, err := os.ReadFile(b.KeyFile)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo da chave privada:", err)
		return false, err
	}

	// Limpeza e verificação do conteúdo da chave privada
	pemKey := strings.TrimSpace(string(keyFileContent))
	if !strings.Contains(pemKey, "-----BEGIN PRIVATE KEY-----") || !strings.Contains(pemKey, "-----END PRIVATE KEY-----") {
		fmt.Println("A chave privada não está no formato PEM esperado.")
		return false, fmt.Errorf("a chave privada não está no formato PEM esperado")
	}

	// Configuração do cliente OCI Object Storage
	configProvider := common.NewRawConfigurationProvider(
		b.Tenancy,
		b.User,
		b.Region,
		b.KeyFingerprint,
		pemKey,
		nil,
	)

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Erro ao criar cliente OCI Object Storage: %v", err)
		panic(err)
	}

	prefix := b.BucketPrefix

	// Listar objetos no bucket para remoção
	request := objectstorage.ListObjectsRequest{
		NamespaceName: common.String(b.NamespaceName),
		BucketName:    common.String(b.BucketName),
		Prefix:        common.String(prefix),
	}

	response, err := client.ListObjects(context.Background(), request)
	if err != nil {
		log.Fatalf("Erro ao listar objetos no bucket: %v", err)
		return false, err
	}

	// Remover os objetos listados
	for _, obj := range response.ListObjects.Objects {
		deleteRequest := objectstorage.DeleteObjectRequest{
			NamespaceName: common.String(b.NamespaceName),
			BucketName:    common.String(b.BucketName),
			ObjectName:    common.String(*obj.Name),
		}

		_, err := client.DeleteObject(context.Background(), deleteRequest)
		if err != nil {
			log.Printf("Erro ao remover o objeto %s: %v", *obj.Name, err)
		} else {
			fmt.Printf("Objeto %s removido com sucesso.\n", *obj.Name)

		}
	}
	return true, nil
}
