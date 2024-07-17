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

func (b BucketConfig) ListarBucket() error {

	keyFileContent, err := os.ReadFile(b.KeyFile)

	if err != nil {
		fmt.Println("Erro ao ler o arquivo da chave privada:", err)
		return err
	}

	// Limpeza e verificação do conteúdo da chave privada
	pemKey := strings.TrimSpace(string(keyFileContent))
	//fmt.Println("pemKey:", pemKey)
	if !strings.Contains(pemKey, "-----BEGIN PRIVATE KEY-----") || !strings.Contains(pemKey, "-----END PRIVATE KEY-----") {
		fmt.Println("A chave privada não está no formato PEM esperado.")
		return err
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
	}

	prefix := b.BucketPrefix

	// Listar objetos no bucket
	request := objectstorage.ListObjectsRequest{
		NamespaceName: common.String(b.NamespaceName),
		BucketName:    common.String(b.BucketName),
		Prefix:        common.String(prefix),
	}

	response, err := client.ListObjects(context.Background(), request)
	if err != nil {
		log.Fatalf("Erro ao listar objetos no bucket: %v", err)
		return err
	}

	// Mostrar os nomes dos objetos (imagens)
	fmt.Printf("Imagens no bucket %s:\n", b.BucketName)
	for _, obj := range response.ListObjects.Objects {
		fmt.Println(*obj.Name)
	}
	return nil
}
