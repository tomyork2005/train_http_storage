package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"RESTarchive/internals/models"
)

func startTestContainer(t *testing.T) string {
	ctx := context.Background()

	absPath, err := filepath.Abs("init.sql")
	require.NoError(t, err)

	r, err := os.Open(absPath)
	require.NoError(t, err)

	req := testcontainers.ContainerRequest{
		Image: "postgres:latest",
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            r,
				HostFilePath:      absPath, // will be discarded internally
				ContainerFilePath: "/init.sql",
			},
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	testcontainers.CleanupContainer(t, postgresC)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port())

	return connStr
}

func TestAddUsersAndUploadFiles(t *testing.T) {

	// Arrange
	tests := []struct {
		name        string
		users       []string
		files       []models.FileToAdd
		obligations []string
		testId      int64
	}{
		{
			name:  "add 2 file to user and get 2 alias",
			users: []string{"andreyTest", "vitaliyTest", "maksTest"},
			files: []models.FileToAdd{
				{
					Alias:      "mustBeAndrey1",
					PathToFile: "nice_path_test1",
					UserId:     1,
				},
				{
					Alias:      "mustBeAndrey2",
					PathToFile: "nice_path_test2",
					UserId:     1,
				},
				{
					Alias:      "mustBeVitaliy1",
					PathToFile: "nice_path_test3",
					UserId:     2,
				},
			},
			obligations: []string{
				"mustBeAndrey1", "mustBeAndrey2",
			},
			testId: 1,
		},
	}

	// Act

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			connectString := startTestContainer(t)

			storage, err := NewStorage(connectString)
			require.NoError(t, err)

			for _, user := range test.users {
				err = storage.NewUser(user)
				require.NoError(t, err)
			}

			for _, file := range test.files {
				err = storage.UploadFiles(file)
				require.NoError(t, err)
			}

			// Assert
			aliases, err := storage.TakeAliasesByUserId(test.testId)
			require.NoError(t, err)

			assert.Equal(t, test.obligations, aliases)

			err = storage.CloseStorage()
			require.NoError(t, err)
		})
	}
}
