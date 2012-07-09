package ec2

import (
	"errors"
	"fmt"
	"github.com/timeredbull/tsuru/config"
	"github.com/timeredbull/tsuru/log"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
	"os"
	"os/user"
	"path"
)

var EC2 *ec2.EC2
var pubKey []byte
var Auth *aws.Auth
var Region *aws.Region

func init() {
	getPubKey()
    Conn()
}

func loadData() {
	var err error
	Auth, err = getAuth()
	if err != nil {
		log.Print(err.Error())
	}
	Region, err = getRegion()
	if err != nil {
		log.Print(err.Error())
	}
}

func getAuth() (*aws.Auth, error) {
	auth := new(aws.Auth)
	var err error
	auth.AccessKey, err = config.GetString("aws:access-key")
	if err != nil {
		msg := "Got error while reaading aws:access-key config options, have you set it?\nError is: " + err.Error()
		log.Print(msg)
		return nil, errors.New(msg)
	}
	auth.SecretKey, err = config.GetString("aws:secret-key")
	if err != nil {
		msg := "Got error while reaading aws:secret-key config options, have you set it?\nError is: " + err.Error()
		log.Print(msg)
		return nil, errors.New(msg)
	}
	return auth, nil
}

func getRegion() (*aws.Region, error) {
	endpnt, err := config.GetString("aws:ec2-endpoint")
	if err != nil {
		msg := "Got error while reaading aws:ec2-endpoint config options, have you set it?\nError is: " + err.Error()
		log.Print(msg)
		return nil, errors.New(msg)
	}
	region := &aws.Region{EC2Endpoint: endpnt}
	return region, nil
}

func getPubKey() ([]byte, error) {
	u, err := user.Current()
	if err != nil {
		return []byte{}, err
	}
	files := []string{"id_dsa.pub", "id_rsa.pub", "identity.pub"}
	for i, f := range files {
		p := path.Join(u.HomeDir, ".ssh", f)
		pubKey, err = ioutil.ReadFile(p)
		if err != nil {
			if os.IsNotExist(err) && i != len(files)-1 {
				continue
			}
			return []byte{}, err
		}
		break
	}
	return pubKey, err
}

func Conn() (*ec2.EC2, error) {
	loadData()
	EC2 = ec2.New(*Auth, *Region)
	return EC2, nil
}

func RunInstance(imageId string, userData string) (string, error) {
	ud := []byte(userData)
	cmd := fmt.Sprintf("\necho \"%s\" >> /root/.ssh/authorized_keys", pubKey)
	ud = append(ud, cmd...)
	rInst := &ec2.RunInstances{
		ImageId:  imageId,
		UserData: ud,
		MinCount: 0,
		MaxCount: 0,
	}
	resp, err := EC2.RunInstances(rInst)
	if err != nil {
		return "", err
	}
	return resp.Instances[0].InstanceId, nil
}
