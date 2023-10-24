package impl

import (
	_errors "errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"sync"
)

type OssGateway struct {
	Host       string
	Port       string
	BucketName string
	BucketID   string
	Lock       *sync.Mutex
	Agent      *gorequest.SuperAgent
}

func (o *OssGateway) GetOsBucketInfos() ([]map[string]interface{}, error) {
	var result = make([]map[string]interface{}, 0)
	if response, body, errors := o.Agent.Clone().Get(fmt.Sprintf("http://%s:%s/api/ossgateway/v1/objectstorageinfo?app=ad", o.Host, o.Port)).EndStruct(&result); len(errors) > 0 {
		return nil, errors[0]
	} else if response.StatusCode != http.StatusOK {
		return nil, _errors.New(string(body))
	}
	return result, nil
}

func (o *OssGateway) GetSingleUploadInfo(filename string) (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	if response, body, errors := o.Agent.Clone().Get(fmt.Sprintf("http://%s:%s/api/ossgateway/v1/upload/%s/%s?type=query_string&request_method=POST&internal_request=true", o.Host, o.Port, o.BucketID, filename)).EndStruct(&result); len(errors) > 0 {
		return nil, errors[0]
	} else if response.StatusCode != http.StatusOK {
		return nil, _errors.New(string(body))
	}
	return result, nil
}

func (o *OssGateway) SingleUpload(content []byte, key string) error {
	//初始化bucketID
	if o.BucketID == "" {
		err := o.initBucketIDByName(o.BucketName)
		if err != nil {
			return err
		}
	}
	//获取上传url
	uploadInfo, err := o.GetSingleUploadInfo(key + ".py")
	if err != nil {
		return err
	}
	var result = make([]map[string]interface{}, 0)
	//文件上传
	url := uploadInfo["url"].(string)
	if response, body, errors := o.Agent.Clone().Post(fmt.Sprintf(url, o.Host, o.Port)).Type(gorequest.TypeMultipart).Send(uploadInfo["form_field"]).SendFile(content, key+".py").EndStruct(&result); len(errors) > 0 {
		return errors[0]
	} else if response.StatusCode != http.StatusOK {
		return _errors.New(string(body))
	}
	return nil
}

func (o *OssGateway) initBucketIDByName(name string) error {
	o.Lock.Lock()
	defer o.Lock.Unlock()
	if o.BucketID != "" {
		return nil
	}
	infos, err := o.GetOsBucketInfos()
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info["name"] == name {
			o.BucketID = fmt.Sprintf("%v", info["id"])
			return nil
		}
	}
	return _errors.New("bucket not found")
}

func (o *OssGateway) GetDownloadInfo(filename string) (string, error) {
	if o.BucketID == "" {
		err := o.initBucketIDByName(o.BucketName)
		if err != nil {
			return "", err
		}
	}
	var result = make(map[string]interface{})
	if response, body, errors := o.Agent.Clone().Get(fmt.Sprintf("http://%s:%s/api/ossgateway/v1/download/%s/%s?type=query_string&internal_request=true", o.Host, o.Port, o.BucketID, filename)).EndStruct(&result); len(errors) > 0 {
		return "", errors[0]
	} else if response.StatusCode != http.StatusOK {
		return "", _errors.New(string(body))
	}
	return result["url"].(string), nil
}

func (o *OssGateway) DeleteFile(filename string) error {
	if o.BucketID == "" {
		err := o.initBucketIDByName(o.BucketName)
		if err != nil {
			return err
		}
	}
	var result = make(map[string]interface{})
	if response, body, errors := o.Agent.Clone().Get(fmt.Sprintf("http://%s:%s/api/ossgateway/v1/delete/%s/%s?internal_request=true", o.Host, o.Port, o.BucketID, filename)).EndStruct(&result); len(errors) > 0 {
		return errors[0]
	} else if response.StatusCode != http.StatusOK {
		return _errors.New(string(body))
	} else {
		headers := result["headers"].(map[string]interface{})
		s := o.Agent.Clone().Delete(result["url"].(string))
		for k, header := range headers {
			s = s.AppendHeader(k, header.(string))
		}
		if response, body, errors = s.EndStruct(&result); len(errors) > 0 {
			return errors[0]
		} else if response.StatusCode != http.StatusOK {
			return _errors.New(string(body))
		}
	}
	return nil
}

func (*OssGateway) DownloadFile(url string) ([]byte, string, error) {
	return []byte{}, "", nil
}
