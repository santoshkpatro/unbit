package cmd

import (
	"context"

	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/santoshkpatro/unbit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	email    string
	fullName string
	password string
)

var addSuperuserCmd = &cobra.Command{
	Use:   "add_superuser",
	Short: "Add Superuser",
	RunE: func(cmd *cobra.Command, args []string) error {
		return addSuperuser(email, fullName, password)
	},
}

func init() {
	addSuperuserCmd.Flags().StringVarP(&email, "email", "e", "", "Email of the superuser")
	addSuperuserCmd.Flags().StringVarP(&fullName, "name", "n", "", "Nameame of the superuser")
	addSuperuserCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the superuser")

	addSuperuserCmd.MarkFlagRequired("email")
	addSuperuserCmd.MarkFlagRequired("name")
	addSuperuserCmd.MarkFlagRequired("password")
}

func addSuperuser(email string, name string, password string) error {
	ctx := context.Background()

	db, err := config.NewPostgresConnection(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	var userExists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err = db.Get(&userExists, query, email)
	if err != nil {
		return err
	}
	if userExists {
		return nil
	}

	salt, _ := utils.GenerateSalt()
	hashedPassword := utils.HashPassword(password, salt)

	insertQuery := `
		INSERT INTO users (id, email, first_name, password_hash, salt, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	_, err = db.Exec(insertQuery, utils.GenerateID("usr"), email, name, hashedPassword, salt, true)
	if err != nil {
		return err
	}

	return nil
}
