package descend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"pault.ag/go/debian/control"
)

func DputFile(client *http.Client, host, archive, fpath string) (error, map[string]string) {
	filename := path.Base(fpath)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://%s/%s/%s", host, archive, filename),
		nil,
	)

	if err != nil {
		return err, map[string]string{}
	}
	fd, err := os.Open(fpath)
	if err != nil {
		return err, map[string]string{}
	}
	req.Body = fd
	resp, err := client.Do(req)
	if err != nil {
		return err, map[string]string{}
	}

	var reply map[string]string
	json.NewDecoder(resp.Body).Decode(&reply)
	return nil, reply
}

func DoPutChanges(client *http.Client, changes *control.Changes, host, archive string) error {
	root := path.Dir(changes.Filename)
	for _, file := range changes.Files {
		err, _ := DputFile(client, host, archive, path.Join(root, file.Filename))
		if err != nil {
			return err
		}
	}
	return nil
}
