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

package main

import (
	"flag"
	"fmt"
	"os"

	"pault.ag/go/config"
	"pault.ag/go/debian/control"
	"pault.ag/go/descend/descend"
)

func Missing(flags *flag.FlagSet, values ...string) {
	for _, value := range values {
		if value != "" {
			continue
		}
		flags.Usage()
		os.Exit(0)
	}
}

func main() {
	conf := descend.Descend{
		Host: "localhost",
		Port: 443,
	}
	flags, err := config.LoadFlags("descend", &conf)
	if err != nil {
		panic(err)
	}
	flags.Parse(os.Args[1:])
	Missing(flags, conf.CaCert, conf.Cert, conf.Key)

	client, err := descend.NewClient(conf.CaCert, conf.Cert, conf.Key)
	if err != nil {
		panic(err)
	}

	for _, changesPath := range flags.Args() {
		changes, err := control.ParseChangesFile(changesPath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Pushing %s\n", changesPath)
		err = descend.DoPutChanges(
			client, changes,
			fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			conf.Archive,
		)
		if err != nil {
			panic(err)
		}
	}
}

// vim: foldmethod=marker
