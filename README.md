# storj-nextcloud

## Initial Set-up
Make sure your `PATH` includes the `$GOPATH/bin` directory, so that your commands can be easily used [Refer: Install the Go Tools](https://golang.org/doc/install):
```
$ export PATH=$PATH:$GOPATH/bin
```

Install dependencies.

```
$ go get
```

## Build ONCE

```
$ go build storj-nextcloud/storj-nextcloud.go
```

## Set-up Files
* Create a `nextcloud_property.json` file with following contents about a NextCloud instance:
    * url :- Login url of the user nextcloud account of the corresponding service provider
    * username :- User ID of the nextcloud account
    * password :- Password of the nextcloud account
```json
    {
        "url": "https://my_nextcloud.com/", //URL should be uptil the third slash only
        "username": "username",
        "password": "password"
  }
```

* Create a `storj_config.json` file with Storj network's configuration information in JSON format:
    * apiKey :- API key created in Storj satellite gui
    * satelliteURL :- Storj Satellite URL
    * encryptionPassphrase :- Storj Encryption Passphrase.
    * bucketName :- Split file into given size before uploading.
    * uploadPath :- Path on Storj Bucket to store data (optional) or "/"
    * serializedScope:- Serialized Scope Key shared while uploading data used to access bucket without API key
    * disallowReads:- Set true to create serialized scope key with restricted read access
    * disallowWrites:- Set true to create serialized scope key with restricted write access
    * disallowDeletes:- Set true to create serialized scope key with restricted delete access
```json
    {

        "apiKey": "change-me-to-the-api-key-created-in-satellite-gui",
        "satelliteURL": "us-central-1.tardigrade.io:7777",
        "bucketName": "change-me-to-desired-bucket-name",
        "uploadPath": "optionalpath/requiredfilename ",
        "encryptionPassphrase": "you'll never guess this",
        "serializedScope": "change-me-to-the-api-key-created-in-encryption-access-apiKey",
        "disallowReads": "true/false-to-disallow-reads",
        "disallowWrites": "true/false-to-disallow-writes",
        "disallowDeletes": "true/false-to-disallow-deletes"
    }
```

* Store both these files in a `config` folder that further is to be stored along with the source file. Filename command-line arguments are optional. Default locations are used.


## Run the command-line tool

**NOTE**: The following commands operate in a Linux system

* Get help
```
    $ ./storj-nextcloud -h
```

* Check version
```
    $ ./storj-nextcloud -v
```

* Connect and transfer file(s)/folder(s) from a desired NextCloud account to desired Storj Bucket
    * **Note** : Path should be `/` for full Back-Up and filepath/folderpath should be in the form `/path/to/file_or_folder`
```
    $ ./storj-nextcloud.go store ./config/nextcloud_property.json ./config/storj_config.json key filePath/folderPath
```

* Read data in `debug` mode from desired NextCloud instance and upload it to given Storj network bucket and download locally to verify successful upload.
```
    $ ./storj-nextcloud.go store debug ./config/nextcloud_property.json ./config/storj_config.json key filePath/folderPath
```

* Read NextCloud instance property from a desired JSON file and display all the contents in the NextCloud account.
    * **Note** : filename arguments are optional. Default locations are used.
```
    $ ./storj-nextcloud.go parse
```

* Read NextCloud instance property in `debug` mode from a desired JSON file and display all the contents in the NextCloud account.
    * **Note** : filename arguments are optional. Default locations are used.
```
    $ ./storj_nextcloud.go parse debug
```

* Read and parse Storj network's configuration, in JSON format, from a desired file and upload a sample object
    * **Note** : filename arguments are optional. Default locations are used.
```
    $ ./storj_nextcloud.go test
```

* Read and parse Storj network's configuration, in JSON format, from a desired file and upload a sample object in `debug` mode to verify successful upload.
    * **Note** : filename arguments are optional. Default locations are used.
```
    $ ./storj_nextcloud.go test debug
```
