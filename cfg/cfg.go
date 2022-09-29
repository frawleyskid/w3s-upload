package cfg

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	ErrMissingIpfsAPI            = errors.New("in `upload` mode, you must supply the multiaddr of a running IPFS node")
	ErrNotUploadMode             = errors.New("some of the given parameters are illegal when not in `upload` mode")
	ErrMissingUploadParams       = errors.New("in `upload` mode, you must supply parameters to configure where to upload files to")
	ErrMissingPinParams          = errors.New("to upload to a pinning service, you must provide both the endpoint URL and your secret token")
	ErrOverloadedPinningServices = errors.New("you cannot upload to both Web3.Storage and a pinning service at the same time")
	ErrInvalidOutput             = errors.New("this is not a valid output type; check the README")
	ErrInvalidMode               = errors.New("this is not a valid mode; check the README")
)

type OperatingMode string

const (
	ModeValidate OperatingMode = "validate"
	ModeHistory  OperatingMode = "history"
	ModeUpload   OperatingMode = "upload"
)

var allModes = map[OperatingMode]struct{}{
	ModeValidate: {},
	ModeHistory:  {},
	ModeUpload:   {},
}

type OutputType string

const (
	OutputArtifacts OutputType = "artifacts"
	OutputCids      OutputType = "cids"
	OutputRootCid   OutputType = "root-cid"
	OutputSummary   OutputType = ""
)

var allOutputs = map[OutputType]struct{}{
	OutputArtifacts: {},
	OutputCids:      {},
	OutputRootCid:   {},
	OutputSummary:   {},
}

type UploadDestination string

const (
	UploadW3S  UploadDestination = "w3s"
	UploadPin  UploadDestination = "pin"
	UploadNone UploadDestination = "none"
)

const (
	DefaultMode = ModeValidate
	DefaultPath = "artifacts/"
)

func init() {
	viper.SetDefault("mode", string(DefaultMode))
	viper.SetDefault("path", string(DefaultPath))

	if err := viper.BindEnv("repo", "GITHUB_WORKSPACE"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("mode", "INPUT_MODE"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("path", "INPUT_PATH"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("w3s-token", "INPUT_W3S-TOKEN"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("ipfs-api", "INPUT_IPFS-API"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("pin-endpoint", "INPUT_PIN-ENDPOINT"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("pin-token", "INPUT_PIN-TOKEN"); err != nil {
		panic(err)
	}

	if err := viper.BindEnv("dry-run", "INPUT_DRY-RUN"); err != nil {
		panic(err)
	}
}

func Repo() string {
	return viper.GetString("repo")
}

func Mode() OperatingMode {
	return OperatingMode(viper.GetString("mode"))
}

func Path() string {
	return viper.GetString("path")
}

func W3SToken() string {
	return viper.GetString("w3s-token")
}

func IpfsAPI() string {
	return viper.GetString("ipfs-api")
}

func PinEndpoint() string {
	return viper.GetString("pin-endpoint")
}

func PinToken() string {
	return viper.GetString("pin-token")
}

func Output() OutputType {
	return OutputType(viper.GetString("output"))
}

func Action() bool {
	return viper.GetBool("action")
}

func DryRun() bool {
	return viper.GetBool("dry-run")
}

func Destination() UploadDestination {
	if viper.IsSet("w3s-token") {
		return UploadW3S
	}

	if viper.IsSet("pin-endpoint") && viper.IsSet("pin-token") {
		return UploadPin
	}

	return UploadNone
}

func ValidateParams() error {
	if _, isValid := allOutputs[Output()]; !isValid {
		return ErrInvalidOutput
	}

	if _, isValid := allModes[Mode()]; !isValid {
		return ErrInvalidMode
	}

	isIpfsAPI := viper.IsSet("ipfs-api")
	isW3sToken := viper.IsSet("w3s-token")
	isPinEndpoint := viper.IsSet("pin-endpoint")
	isPinToken := viper.IsSet("pin-token")

	if Mode() == ModeUpload {
		switch {
		case !isIpfsAPI:
			return ErrMissingIpfsAPI
		case !(isW3sToken || isPinEndpoint || isPinToken):
			return ErrMissingUploadParams
		case isW3sToken && (isPinEndpoint || isPinToken):
			return ErrOverloadedPinningServices
		case isPinEndpoint != isPinToken:
			return ErrMissingPinParams
		}
	} else if isIpfsAPI || isW3sToken || isPinEndpoint || isPinToken {
		return ErrNotUploadMode
	}

	return nil
}
