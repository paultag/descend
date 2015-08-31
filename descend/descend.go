package descend

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"pault.ag/go/debian/control"
)

func DputFile(client *http.Client, host, archive, fpath string) error {
	filename := path.Base(fpath)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://%s/%s/%s", host, archive, filename),
		nil,
	)

	if err != nil {
		return err
	}
	fd, err := os.Open(fpath)
	if err != nil {
		return err
	}
	req.Body = fd
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func DoPutChanges(client *http.Client, changes *control.Changes, host, archive string) error {
	root := path.Dir(changes.Filename)
	for _, file := range changes.Files {
		err := DputFile(client, host, archive, path.Join(root, file.Filename))
		if err != nil {
			return err
		}
	}
	return nil
}
