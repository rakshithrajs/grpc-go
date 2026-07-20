package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rakshithrajs/cloud/services/files/internal/config"

	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	token       string
	filesAddr   string
	filesClient filespb.FilesClient
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	flag.StringVar(&token, "token", "", "JWT token for authentication (or FILES_TOKEN env)")
	flag.StringVar(&filesAddr, "addr", cfg.GRPCAddress, "files gRPC server address")
	flag.Parse()

	if token == "" {
		fmt.Printf("missing authentication token: provide --token flag or set FILES_TOKEN env var\n")
		os.Exit(1)
	}

	conn, err := grpc.NewClient(filesAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("failed to create client connection: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	filesClient = filespb.NewFilesClient(conn)

	if err := runInteractive(); err != nil {
		fmt.Printf("client error: %v\n", err)
		os.Exit(1)
	}
}

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func withAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}

func prompt(label string) string {
	fmt.Printf("%s: ", label)
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func runInteractive() error {
	for {
		fmt.Println()
		fmt.Println("Files service client")
		fmt.Println("1. Upload a file")
		fmt.Println("2. List files")
		fmt.Println("3. Download a file")
		fmt.Println("4. Rename a file")
		fmt.Println("5. Delete a file")
		fmt.Println("6. Quit")
		fmt.Println()

		choice := prompt("Choose an option")

		var runErr error
		var action string
		switch choice {
		case "1":
			action = "upload"
			runErr = promptAndUploadFile()
		case "2":
			action = "list"
			runErr = runListFiles()
		case "3":
			action = "download"
			runErr = promptAndDownloadFile()
		case "4":
			action = "rename"
			runErr = promptAndRenameFile()
		case "5":
			action = "delete"
			runErr = promptAndDeleteFile()
		case "6":
			fmt.Println("exiting")
			return nil
		default:
			fmt.Println("invalid option, please choose 1-6")
			continue
		}

		if runErr != nil {
			fmt.Printf("%s: %v\n", action, runErr)
		}
	}
}

func promptAndUploadFile() error {
	filePath := prompt("File path")
	if filePath == "" {
		return fmt.Errorf("file path is required")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = withAuth(ctx)

	resp, err := filesClient.UploadFile(ctx, &filespb.UploadFileRequest{
		FileName: filepath.Base(filePath),
		Content:  content,
	})
	if err != nil {
		return err
	}

	fmt.Printf("file uploaded: %s (%d bytes, %s)\n",
		resp.GetFile().GetFileName(),
		resp.GetFile().GetFileSize(),
		resp.GetFile().GetMimeType(),
	)
	fmt.Printf("file ID: %s\n", resp.GetFile().GetID())
	return nil
}

func runListFiles() error {
	ctx, cancel := withTimeout()
	defer cancel()
	ctx = withAuth(ctx)

	resp, err := filesClient.ListFiles(ctx, &filespb.ListFilesRequest{})
	if err != nil {
		return err
	}

	if len(resp.GetFile()) == 0 {
		fmt.Println("no files found")
		return nil
	}

	fmt.Println("Files")
	for _, file := range resp.GetFile() {
		fmt.Printf("%s | %s | %d bytes | %s\n",
			file.GetID(),
			file.GetFileName(),
			file.GetFileSize(),
			file.GetMimeType(),
		)
	}
	return nil
}

func promptAndDownloadFile() error {
	fileID := prompt("File ID")
	if fileID == "" {
		return fmt.Errorf("file ID is required")
	}

	outputPath := prompt("Output path")
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = withAuth(ctx)

	resp, err := filesClient.DownloadFile(ctx, &filespb.DownloadFileRequest{FileID: fileID})
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(outputPath, resp.GetFileName()), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(resp.GetContent()); err != nil {
		return fmt.Errorf("failed to write downloaded file: %w", err)
	}

	fmt.Printf("file downloaded: %s -> %s\n", resp.GetFileName(), outputPath)
	return nil
}

func promptAndRenameFile() error {
	fileID := prompt("File ID")
	if fileID == "" {
		return fmt.Errorf("file ID is required")
	}

	newName := prompt("New file name")
	if newName == "" {
		return fmt.Errorf("new file name is required")
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = withAuth(ctx)

	_, err := filesClient.RenameFile(ctx, &filespb.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	})
	if err != nil {
		return err
	}

	fmt.Printf("file renamed: %s -> %s\n", fileID, newName)
	return nil
}

func promptAndDeleteFile() error {
	fileID := prompt("File ID")
	if fileID == "" {
		return fmt.Errorf("file ID is required")
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = withAuth(ctx)

	_, err := filesClient.DeleteFile(ctx, &filespb.DeleteFileRequest{FileID: fileID})
	if err != nil {
		return err
	}

	fmt.Printf("file deleted: %s\n", fileID)
	return nil
}
