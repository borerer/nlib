package utils

import "github.com/mitchellh/mapstructure"

func DecodeStruct(input interface{}, res interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		Result:  res,
		TagName: "json",
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}
	if err := decoder.Decode(input); err != nil {
		return err
	}
	return nil
}
