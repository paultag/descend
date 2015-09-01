/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

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
	putPath := path.Join(archive, filename)
	url := fmt.Sprintf("https://%s/%s", host, putPath)
	req, err := http.NewRequest("PUT", url, nil)

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

// vim: foldmethod=marker
