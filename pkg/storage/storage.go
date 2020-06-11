/*
Copyright © 2020 Arnoud Kleinloog

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package storage

import (
	"encoding/json"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/spf13/afero"
)

func Get(key string) (interface{}, bool, error) {

	exists, err := afero.Exists(app.Fs, key)
	if err != nil {
		app.Log.Error(err, "Error occurred while checking if content exists")
		return nil, false, err
	}

	if exists {

		isDir, err := afero.IsDir(app.Fs, key)
		if err != nil {
			app.Log.Error(err, "Error occurred while checking if directory")
			return nil, false, err
		}
		if isDir {
			return nil, false, nil
		}

		bytes, err := afero.ReadFile(app.Fs, key)
		if err != nil {
			app.Log.Error(err, "Error occurred while reading content")
			return nil, true, err
		}

		var content interface{}
		err = json.Unmarshal(bytes, &content)
		if err != nil {
			app.Log.Error(err, "Error occurred while unmarshalling content from JSON")
			return nil, true, err
		}

		return content, true, nil
	}

	return nil, false, nil
}

func GetCollection(key string) (map[string]interface{}, error) {

	var itemsInCollection = make(map[string]interface{})

	exists, err := afero.DirExists(app.Fs, key)
	if err != nil {
		app.Log.Error(err, "Error occurred while checking if directory exists")
		return nil, err
	}

	if exists {
		files, err := afero.ReadDir(app.Fs, key)
		if err != nil {
			app.Log.Error(err, "Error occurred while retrieving files in directory")
			return nil, err
		}

		for index := range files {
			fileInfo := files[index]
			if !fileInfo.IsDir() {
				content, exists, err := Get(key + "/" + fileInfo.Name())

				if err != nil {
					app.Log.Error(err, "Error occurred while retrieving individual file in directory")
					return nil, err
				}

				if exists {
					itemsInCollection[fileInfo.Name()] = content
				}
			}
		}
	}

	return itemsInCollection, nil
}

func Store(key string, content interface{}) error {

	bytes, err := json.Marshal(content)
	if err != nil {
		app.Log.Error(err, "Error marshalling content to JSON")
		return err
	}

	err = afero.WriteFile(app.Fs, key, bytes, 0644)
	if err != nil {
		app.Log.Error(err, "Error occurred while storing content")
		return err
	}

	return nil
}

func Remove(key string) (bool, error) {

	exists, err := afero.Exists(app.Fs, key)
	if err != nil {
		app.Log.Error(err, "Error occurred while checking if content exists")
		return false, err
	}
	if exists {
		err = app.Fs.Remove(key)
		if err != nil {
			app.Log.Error(err, "Error occurred while removing content")
			return false, err
		}
	}
	return exists, nil
}
