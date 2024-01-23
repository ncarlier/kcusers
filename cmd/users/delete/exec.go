package deleteuser

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ncarlier/kcusers/pkg/keycloak"
	"github.com/ncarlier/kcusers/pkg/uuid"
)

const unableToDeleteUsers = "unable to delete users"

func exec(client *keycloak.Client, params execParams) error {
	file, err := os.Open(params.filename)
	if err != nil {
		return fmt.Errorf("%s: %w", unableToDeleteUsers, err)
	}
	defer file.Close()
	slog.Info("deleting users...", "filename", params.filename)

	sem := make(chan struct{}, params.concurent)
	defer close(sem)

	wg := sync.WaitGroup{}

	var total, deleted, errors int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		uid, ok := uuid.GetUUIDPrefix(line)
		if !ok {
			slog.Error("invalid line: skiping")
			continue
		}

		sem <- struct{}{} // Acquire a semaphore slot
		slog.Debug("deleting user", "uid", uid)
		total++
		wg.Add(1)
		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()
			if params.dryRun {
				slog.Info("user dry run deletion", "uid", uid)
				return
			}
			now := time.Now()
			if err := deleteUser(client, uid); err != nil {
				slog.Error(unableToDeleteUsers, "uid", uid, "error", err)
				errors++
				return
			}
			deleted++
			slog.Info("user deleted", "uid", uid, "took", time.Since(now).Milliseconds())
		}()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%s: %w", unableToDeleteUsers, err)
	}

	wg.Wait()

	slog.Info("users deleted", "total", total, "deleted", deleted, "errors", errors)

	return nil
}

func deleteUser(client *keycloak.Client, uid string) error {
	endpoint := fmt.Sprintf("%s/users/%s", client.GetAdminBaseURL(), uid)
	req, err := http.NewRequest("DELETE", endpoint, http.NoBody)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("invalid response: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
