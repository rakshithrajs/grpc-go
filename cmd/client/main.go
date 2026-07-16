package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	authpb "cloud/gen/auth/v1"
	filespb "cloud/gen/files/v1"
	userspb "cloud/gen/users/v1"
	"cloud/internal/config"

	"github.com/charmbracelet/huh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var token string
var authClient authpb.AuthClient
var userClient userspb.UsersClient
var fileClient filespb.FilesClient

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		printError("failed to load config: %v", err)
		os.Exit(1)
	}

	conn, err := grpc.NewClient(cfg.GRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		printError("failed to create client connection: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	authClient = authpb.NewAuthClient(conn)
	userClient = userspb.NewUsersClient(conn)
	fileClient = filespb.NewFilesClient(conn)

	if err := runInteractive(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			printInfo("exiting")
			return
		}
		printError("client error: %v", err)
		os.Exit(1)
	}
}

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func runInteractive() error {
	for {
		var action string
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Title("What would you like to do?").Options(
					huh.NewOption("Register a new user", "register"),
					huh.NewOption("Log in an existing user", "login"),
					huh.NewOption("Get user profile", "getProfile"),
					huh.NewOption("Update user profile", "updateUserProfile"),
					huh.NewOption("Upload a file", "uploadFile"),
					huh.NewOption("List files", "listFiles"),
					huh.NewOption("Download a file", "downloadFile"),
					huh.NewOption("Rename a file", "renameFile"),
					huh.NewOption("Delete a file", "deleteFile"),
					huh.NewOption("Quit", "quit"),
				).Value(&action),
			),
		).Run(); err != nil {
			return err
		}

		if action == "quit" {
			return nil
		}

		var runErr error
		switch action {
		case "register":
			runErr = promptAndRegister()
		case "login":
			runErr = promptAndLogin()
		case "getProfile":
			runErr = runGetProfile()
		case "updateUserProfile":
			runErr = promptAndUpdateUserProfile()
		case "uploadFile":
			runErr = promptAndUploadFile()
		case "listFiles":
			runErr = runListFiles()
		case "downloadFile":
			runErr = promptAndDownloadFile()
		case "renameFile":
			runErr = promptAndRenameFile()
		case "deleteFile":
			runErr = promptAndDeleteFile()
		default:
			runErr = fmt.Errorf("unknown action: %s", action)
		}

		if runErr != nil {
			printError("%s: %v", action, runErr)
		}
	}
}

func promptAndRegister() error {
	var name, email, password, phone string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Placeholder("Alice Smith").
				Value(&name).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Email").
				Placeholder("alice@example.com").
				Value(&email).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Password").
				Placeholder("Secret@123").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Phone").
				Placeholder("1234567890").
				Value(&phone).
				Validate(huh.ValidateNotEmpty()),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	ctx, cancel := withTimeout()
	defer cancel()
	return runRegister(ctx, name, email, password, phone)
}

func promptAndUpdateUserProfile() error {
	var name, email, phone, password string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Placeholder("Alice Smith").
				Value(&name),
			huh.NewInput().
				Title("Email").
				Placeholder("alice@example.com").
				Value(&email),
			huh.NewInput().
				Title("Phone").
				Placeholder("1234567890").
				Value(&phone),
			huh.NewInput().
				Title("Password").
				Placeholder("Secret@123").
				EchoMode(huh.EchoModePassword).
				Value(&password),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if name == "" && email == "" && phone == "" && password == "" {
		return errors.New("at least one field must be provided")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	return runUpdateUserProfile(ctx, name, email, phone, password)
}

func promptAndLogin() error {
	var email, password string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Placeholder("alice@example.com").
				Value(&email).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Password").
				Placeholder("Secret@123").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Validate(huh.ValidateNotEmpty()),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	ctx, cancel := withTimeout()
	defer cancel()
	return runLogin(ctx, email, password)
}

func runRegister(ctx context.Context, name, email, password, phone string) error {
	resp, err := authClient.RegisterUser(ctx, &authpb.RegisterUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		Phone:    phone,
	})
	if err != nil {
		return err
	}
	printSuccess("register successful: %s (%s)", resp.GetUser().GetName(), resp.GetUser().GetEmail())
	return nil
}

func runLogin(ctx context.Context, email, password string) error {
	resp, err := authClient.LoginUser(ctx, &authpb.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}
	token = resp.GetToken()
	printSuccess("login successful")
	return nil
}

func runGetProfile() error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := userClient.GetUserProfile(ctx, &userspb.GetUserProfileRequest{})
	if err != nil {
		return err
	}
	user := resp.GetUser()
	printHeader("User Profile")
	printInfo("ID:    %s", user.GetId())
	printInfo("Name:  %s", user.GetName())
	printInfo("Email: %s", user.GetEmail())
	printInfo("Phone: %s", user.GetPhone())
	return nil
}

func runUpdateUserProfile(ctx context.Context, name, email, phone, password string) error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	_, err := userClient.UpdateUserProfile(ctx, &userspb.UpdateUserProfileRequest{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
	})
	if err != nil {
		return err
	}
	printSuccess("profile updated")
	return nil
}

func promptAndUploadFile() error {
	var filePath string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("File path").
				Placeholder("paste file path here on your local systen").
				Value(&filePath).
				Validate(huh.ValidateNotEmpty()),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	return runUploadFile(filePath)
}

func runUploadFile(filePath string) error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := fileClient.UploadFile(ctx, &filespb.UploadFileRequest{
		FileName: filepath.Base(filePath),
		Content:  content,
	})
	if err != nil {
		return err
	}

	printSuccess("file uploaded: %s (%d bytes, %s)",
		resp.GetFile().GetFileName(),
		resp.GetFile().GetFileSize(),
		resp.GetFile().GetMimeType(),
	)
	return nil
}

func runListFiles() error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := fileClient.ListFiles(ctx, &filespb.ListFilesRequest{})
	if err != nil {
		return err
	}

	if len(resp.GetFile()) == 0 {
		printInfo("no files found")
		return nil
	}

	printHeader("Files")
	for _, file := range resp.GetFile() {
		printInfo("%s | %s | %d bytes | %s",
			file.GetID(),
			file.GetFileName(),
			file.GetFileSize(),
			file.GetMimeType(),
		)
	}
	return nil
}

func promptAndDownloadFile() error {
	var fileID, outputPath string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("File ID").
				Placeholder("paste the ID from ListFiles").
				Value(&fileID).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Output path").
				Placeholder("where to save the downloaded file").
				Value(&outputPath).
				Validate(huh.ValidateNotEmpty()),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	return runDownloadFile(fileID, outputPath)
}

func runDownloadFile(fileID, outputPath string) error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := fileClient.DownloadFile(ctx, &filespb.DownloadFileRequest{FileID: fileID})
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, resp.GetContent(), 0o644); err != nil {
		return fmt.Errorf("failed to write downloaded file: %w", err)
	}

	printSuccess("file downloaded: %s → %s", resp.GetFileName(), outputPath)
	return nil
}

func promptAndRenameFile() error {
	var fileID, newName string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("File ID").
				Placeholder("paste the ID from ListFiles").
				Value(&fileID).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("New file name").
				Placeholder("new name without path separators").
				Value(&newName).
				Validate(huh.ValidateNotEmpty()),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	return runRenameFile(fileID, newName)
}

func runRenameFile(fileID, newName string) error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	_, err := fileClient.RenameFile(ctx, &filespb.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	})
	if err != nil {
		return err
	}

	printSuccess("file renamed: %s → %s", fileID, newName)
	return nil
}

func promptAndDeleteFile() error {
	var fileID string
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("File ID").
				Placeholder("paste the ID from ListFiles").
				Value(&fileID).
				Validate(huh.ValidateNotEmpty()),
			huh.NewConfirm().
				Title("Are you sure you want to delete this file?").
				Value(&confirm).
				Affirmative("Yes").
				Negative("No"),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if !confirm {
		printInfo("delete cancelled")
		return nil
	}
	return runDeleteFile(fileID)
}

func runDeleteFile(fileID string) error {
	if token == "" {
		return errors.New("not logged in: please log in first")
	}

	ctx, cancel := withTimeout()
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
	_, err := fileClient.DeleteFile(ctx, &filespb.DeleteFileRequest{FileID: fileID})
	if err != nil {
		return err
	}

	printSuccess("file deleted: %s", fileID)
	return nil
}
