package files

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.cto.ai/provision/internal/awsrds"

	"git.cto.ai/provision/internal/setup"

	"git.cto.ai/provision/internal/logger"
	ctoai "github.com/cto-ai/sdk-go"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func EBRepoFileSetup(ux *ctoai.Ux, githubRepoDetails setup.GithubRepoDetails, rdsBool bool, rdsDetails awsrds.RDSDetails) (string, error) {
	githubDownloadLink, err := getDownloadLink(githubRepoDetails)
	if err != nil {
		return "", err
	}

	err = download(ux, fmt.Sprintf("%s.zip", githubRepoDetails.Repo), githubDownloadLink)
	if err != nil {
		return "", err
	}

	unzippedRepo, err := unzip(fmt.Sprintf("%s.zip", githubRepoDetails.Repo))
	if err != nil {
		return unzippedRepo, err
	}

	if rdsBool {
		content := fmt.Sprintf(`option_settings:
  - option_name: RDS_HOSTNAME
    value: %s 
  - option_name: RDS_USERNAME 
    value: %s 
  - option_name: RDS_PASSWORD 
    value: %s 
  - option_name: RDS_PORT 
    value: %s 
  - option_name: RDS_DB_NAME 
    value: %s`, rdsDetails.Host, rdsDetails.Username, rdsDetails.Password, rdsDetails.Port, rdsDetails.DBName)

		err := createEBExtentions(content, unzippedRepo, "rds_env")
		if err != nil {
			return "", err
		}
	}

	err = rezip(unzippedRepo)
	if err != nil {
		return unzippedRepo, err
	}

	return unzippedRepo, nil
}

func getDownloadLink(githubRepoDetails setup.GithubRepoDetails) (string, error) {
	if githubRepoDetails.Token != "public" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubRepoDetails.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		s := github.RepositoryContentGetOptions{}

		archiveLink, _, err := client.Repositories.GetArchiveLink(ctx, githubRepoDetails.Username, githubRepoDetails.Repo, github.Zipball, &s, false)
		if err != nil {
			return "", err
		}

		return archiveLink.String(), nil
	}

	return fmt.Sprintf("http://github.com/%s/%s/zipball/master", githubRepoDetails.Username, githubRepoDetails.Repo), nil
}

func download(ux *ctoai.Ux, filepath string, url string) error {
	logger.LogSlack(ux, "🔄 Downloading repository files...")

	resp, err := http.Get(url)
	if err != nil {
		logger.LogSlackError(ux, err)
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		logger.LogSlackError(ux, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		logger.LogSlackError(ux, err)
		return err
	}

	logger.LogSlack(ux, "✅ Download complete.")
	return err
}

func unzip(src string) (string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := f.Name

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}

	unzippedRepo := strings.Replace(filenames[0], "/", "", -1)

	return unzippedRepo, nil
}

func rezip(dir string) error {
	zipCmd := fmt.Sprintf("zip ../%s.zip -r * .[^.]*", dir)
	cmd := exec.Command("sh", "-c", zipCmd)
	cmd.Dir = fmt.Sprintf("./%s", dir)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createEBExtentions(content, unzippedRepo, extName string) error {
	dirs := fmt.Sprintf("./%s/.ebextensions/", unzippedRepo)
	os.MkdirAll(dirs, os.ModePerm)
	f, err := os.Create(dirs + fmt.Sprintf("%s.config", extName))
	if err != nil {
		return err
	}

	dbSettings := content
	d2 := []byte(dbSettings)
	_, err = f.Write(d2)
	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		f.Close()
		return err
	}

	return nil
}
