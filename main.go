package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Xato: .env faylni yuklashda muammo:", err)
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Fatalln("Xato: .env faylda MINIO_ENDPOINT, MINIO_ACCESS_KEY yoki MINIO_SECRET_KEY topilmadi")
	}


	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Xato:", err)
	}

	log.Println("MinIO ga muvaffaqiyatli ulandik!")

	ctx := context.Background()

	var bucketName string
	fmt.Print("Bucket nomini kiriting: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	bucketName = scanner.Text()
	
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Println("Bucket yaratishda xatolik:", err)
	} else {
		log.Println("Bucket muvaffaqiyatli yaratildi:", bucketName)
	}

	fileContent := "Salom bu test file"
	reader := strings.NewReader(fileContent)
	objectName := "text.txt"
	_, err = minioClient.PutObject(ctx, bucketName, objectName, reader, int64(len(fileContent)), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		log.Fatalln("Fayil yuklash xatolik:", err)
	}

	log.Println("Fayil yuklandi:", objectName)
}

