package auth_storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"utils"
)

var (
	keys         *Keys
	keys_lock    sync.Mutex
	key_filepath string
)

func InitKey(key_filepath string) {
	// 使用默认key
	json.Unmarshal(defaultKey, keys)
	if key_filepath == "" {
		fmt.Printf("auth key path is empty, will use default key.\n")
	} else {
		err := loadKeyFromFile(key_filepath)
		if err != nil {
			fmt.Printf("load keys failed: %v\n", err)
			return
		}
	}
}

func loadKeyFromFile(key_filepath string) error {
	keys_lock.Lock()
	defer keys_lock.Unlock()
	if utils.PathExist(key_filepath) == false {
		return fmt.Errorf("%s not exist\n", key_filepath)
	}
	content, err := ioutil.ReadFile(key_filepath)
	if err != nil {
		return fmt.Errorf(" read file %s content error %v", key_filepath, err)
	}
	keys = &Keys{}
	err = json.Unmarshal(content, keys)
	if err != nil {
		return fmt.Errorf(" unmarshal file %s content error %v", key_filepath, err)
	}
	return nil
}

func Query(outKey string) (*Key, error) {
	keys_lock.Lock()
	defer keys_lock.Unlock()
	if keys == nil {
		return nil, fmt.Errorf("key is null")
	}
	if len(keys.Keys) == 0 {
		return nil, fmt.Errorf("key length is 0")
	}
	for _, k := range keys.Keys {
		if k.Access == outKey {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("can not find key")
}

var defaultKey = []byte(`
{
  "version": "1.0.0",
  "keys": [
    {
      "access": "1247fbc0e373637ff9fb08f0c81722f3",
      "secret": "d8e8fca2dc0f896fd7cb4cb0031ba249",
      "product_list": [
        "fake_A04"
      ]
    },
    {
      "access": "6c05c3c9904e4d499917cf26278b3741",
      "secret": "ccccf8ef-82fd-421b-83e1-4d5cf47b575a",
      "product_list": [
        "A01"
      ]
    },
    {
      "access": "be6a470094354f2f8c09414cb1b8d73f",
      "secret": "220ad5e6-560e-4ba9-84d3-46569f851e1a",
      "product_list": [
        "A02"
      ]
    },
    {
      "access": "14aaa4312a3a42f09e723f580ce8019c",
      "secret": "2f725236-a29b-4ffc-82e7-71a65fe0345b",
      "product_list": [
        "A03"
      ]
    },
    {
      "access": "c2533fd4d6944ebda83244ca61c9de83",
      "secret": "f8a0e6bb-840d-4b97-ad05-edd34da853cd",
      "product_list": [
        "A04"
      ]
    },
    {
      "access": "74bbf69280a44481b4da0bc2c0035b3b",
      "secret": "56cb54cc-74de-4d13-b9e2-98d511115e47",
      "product_list": [
        "A05"
      ]
    },
    {
      "access": "6488127ad4f346e0b09ee2459af34301",
      "secret": "e6db2ad4-bfc5-4a1f-903b-052dc3bec6df",
      "product_list": [
        "A06"
      ]
    },
    {
      "access": "bc58d50449524a689f635636451cda2c",
      "secret": "5df784a1-94cf-4a18-96bf-bb6388af7855",
      "product_list": [
        "B01"
      ]
    },
    {
      "access": "229a0417627148fab3cfc96b71b8c185",
      "secret": "0a01917e-fc52-4f23-abd0-01b10a2142e3",
      "product_list": [
        "B05"
      ]
    },
    {
      "access": "04af0ffc48494ea9a7bb597c25f1293d",
      "secret": "ca5d62c9-8f79-456a-a929-c94e4189da5f",
      "product_list": [
        "B06"
      ]
    },
    {
      "access": "a8c094c041f5494681e94443604fe59b",
      "secret": "0610747c-b3d3-44cc-af73-b1cd59836ebb",
      "product_list": [
        "B07"
      ]
    },
    {
      "access": "059f0d43f19343ccb27cd15eee02b37c",
      "secret": "224989e1-67f7-46f7-84e2-6457c420352e",
      "product_list": [
        "C01"
      ]
    },
    {
      "access": "d9383e1002e142ce8824c13140028002",
      "secret": "ec068469-4c7b-43f0-8ebc-7b19eff0422c",
      "product_list": [
        "C02"
      ]
    },
    {
      "access": "810d3c7cb8dc4cdc8092635c31618f91",
      "secret": "3a6d7af4-108e-4294-961b-28572197816a",
      "product_list": [
        "C07"
      ]
    },
    {
      "access": "abad35b421d84739aad387213b87be6f",
      "secret": "31df8864-76bd-4ac6-bdcd-42bd4dbccf00",
      "product_list": [
        "A01",
        "A02",
        "A03",
        "A04",
        "A05",
        "A06",
        "B01",
        "B05",
        "B06",
        "B07",
        "C01",
        "C02",
        "C07"
      ]
    }
  ]
}
`)
