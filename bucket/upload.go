package bucket

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

// BucketConfig struct to hold environment variables
type BucketConfig struct {
	Tenancy        string
	User           string
	KeyFile        string
	KeyFingerprint string
	Region         string
	NamespaceName  string
	BucketName     string
	BucketPrefix   string
}

func LoadConfig() (*BucketConfig, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Map environment variables to Config struct
	config := &BucketConfig{
		Tenancy:        os.Getenv("OCI_TENANCY"),
		User:           os.Getenv("OCI_USER"),
		KeyFile:        os.Getenv("OCI_KEY_FILE"),
		KeyFingerprint: os.Getenv("OCI_KEY_FINGERPRINT"),
		Region:         os.Getenv("OCI_REGION"),
		NamespaceName:  os.Getenv("NAMESPACE_NAME"),
		BucketName:     os.Getenv("BUCKET_NAME"),
		BucketPrefix:   os.Getenv("BUCKET_PREFIX"),
	}

	return config, nil
}

func (b BucketConfig) UploadObject(objectName string) (string, error) {

	prefix := b.BucketPrefix
	newName := strings.Trim(objectName, "/images")

	objCompleteName := prefix + newName

	// Ler conteúdo do arquivo da chave privada
	keyFileContent, err := os.ReadFile(b.KeyFile)

	if err != nil {
		fmt.Println("Erro ao ler o arquivo da chave privada:", err)
		panic(err)
	}

	// Limpeza e verificação do conteúdo da chave privada
	pemKey := strings.TrimSpace(string(keyFileContent))

	if !strings.Contains(pemKey, "-----BEGIN PRIVATE KEY-----") || !strings.Contains(pemKey, "-----END PRIVATE KEY-----") {

		panic("private key is not in expected PEM format")
	}

	// Converter a chave para uma única linha, se necessário

	// Configuração de autenticação com o conteúdo da chave privada

	provider := common.NewRawConfigurationProvider(
		b.Tenancy,
		b.User,
		b.Region,
		b.KeyFingerprint,
		pemKey, // Passar o conteúdo da chave como string
		nil,
	)

	// Configuração do cliente OCI Object Storage
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		panic("Erro ao criar cliente OCI Object Storage")
	}

	// Nome do bucket e caminho do arquivo local
	//bucketName := os.Getenv("BUCKET_NAME")
	filePath := objectName // Substitua pelo caminho real do arquivo

	// Abrir o arquivo local para leitura
	file, err := os.Open(filePath)
	if err != nil {

		return "", fmt.Errorf("Erro ao abrir arquivo local", err)
	}
	defer file.Close()

	// Preparar para fazer upload do arquivo
	// Nome do arquivo no Object Storage

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("Erro ao obter informações do arquivo: %v", err)

	}

	// Criar objeto para upload
	request := objectstorage.PutObjectRequest{
		NamespaceName: common.String(b.NamespaceName),
		BucketName:    common.String(b.BucketName),
		ObjectName:    common.String(objCompleteName),
		PutObjectBody: file,
		ContentLength: common.Int64(fileInfo.Size()), // Definir o tamanho do arquivo
	}

	// Realizar upload do objeto
	_, err = client.PutObject(context.Background(), request)
	if err != nil {

		return "", fmt.Errorf("Erro ao realizar upload do objeto:", err)
	}

	deleteFile(filePath)

	return ("Upload de " + objectName + " realizado com sucesso."), nil

}

func deleteFile(filePath string) error {
	// Tenta deletar o arquivo
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("erro ao deletar o arquivo '%s': %w", filePath, err)
	}
	return nil
}
