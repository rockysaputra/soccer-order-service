package util

import (
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BindFromJson(dest any, filename string, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal : %s", err)
		return err
	}

	return nil
}

func setEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)

	if err != nil {
		logrus.Errorf("failed to unmarshal : %s", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			panic("unsupported type")
		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env : %s", err)
			return err
		}
	}

	return nil
}

func BindFromConsulKv(dest any, endpoint, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endpoint, path)

	if err != nil {
		logrus.Errorf("failed to add remote provider : %s", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote config : %s", err)
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal : %s", err)
		return err
	}

	err = setEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("failed to set env from consul kv : %s", err)
		return err
	}

	return nil
}
