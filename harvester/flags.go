package harvester

import (
	"errors"
	"fmt"

	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
)

const (
	defaultNamespace = "default"

	defaultPort = 30443

	defaultCPU          = 2
	defaultMemorySize   = 2
	defaultDiskSize     = 20
	defaultDiskBus      = "virtio"
	defaultNetworkModel = "virtio"
	networkTypePod      = ""
	networkTypeDHCP     = "dhcp"
	networkTypeStatic   = "static"
)

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_HOST",
			Name:   "harvester-host",
			Usage:  "harvester host",
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_PORT",
			Name:   "harvester-port",
			Usage:  "harvester port",
			Value:  defaultPort,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_USERNAME",
			Name:   "harvester-username",
			Usage:  "harvester username",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_PASSWORD",
			Name:   "harvester-password",
			Usage:  "harvester password",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NAMESPACE",
			Name:   "harvester-namespace",
			Usage:  "harvester namespace",
			Value:  defaultNamespace,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_CPU_COUNT",
			Name:   "harvester-cpu-count",
			Usage:  "number of CPUs for the machine",
			Value:  defaultCPU,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_MEMORY_SIZE",
			Name:   "harvester-memory-size",
			Usage:  "size of memory for machine (in GiB)",
			Value:  defaultMemorySize,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_DISK_SIZE",
			Name:   "harvester-disk-size",
			Usage:  "size of disk for machine (in GiB)",
			Value:  defaultDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_DISK_BUS",
			Name:   "harvester-disk-bus",
			Usage:  "bus of disk for machine",
			Value:  defaultDiskBus,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_IMAGE_NAME",
			Name:   "harvester-image-name",
			Usage:  "harvester image name",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_SSH_USER",
			Name:   "harvester-ssh-user",
			Usage:  "SSH username",
			Value:  drivers.DefaultSSHUser,
		},
		mcnflag.IntFlag{
			EnvVar: "HARVESTER_SSH_PORT",
			Name:   "harvester-ssh-port",
			Usage:  "SSH port",
			Value:  drivers.DefaultSSHPort,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_SSH_PASSWORD",
			Name:   "harvester-ssh-password",
			Usage:  "SSH password",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_KEY_PAIR_NAME",
			Name:   "harvester-key-pair-name",
			Usage:  "harvester key pair name",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_SSH_PRIVATE_KEY_PATH",
			Name:   "harvester-ssh-private-key-path",
			Usage:  "SSH private key path ",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_TYPE",
			Name:   "harvester-network-type",
			Usage:  "harvester network type",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_NAME",
			Name:   "harvester-network-name",
			Usage:  "harvester network name",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_MODEL",
			Name:   "harvester-network-model",
			Usage:  "harvester network model",
			Value:  defaultNetworkModel,
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_INTERFACE",
			Name:   "harvester-network-interface",
			Usage:  "harvester network interface",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_MASK",
			Name:   "harvester-network-mask",
			Usage:  "harvester network mask",
		},
		mcnflag.StringFlag{
			EnvVar: "HARVESTER_NETWORK_GATEWAY",
			Name:   "harvester-network-gateway",
			Usage:  "harvester network gateway",
		},
	}
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Host = flags.String("harvester-host")
	d.Port = flags.Int("harvester-port")
	d.Username = flags.String("harvester-username")
	d.Password = flags.String("harvester-password")
	d.Namespace = flags.String("harvester-namespace")

	d.CPU = flags.Int("harvester-cpu-count")
	d.MemorySize = fmt.Sprintf("%dGi", flags.Int("harvester-memory-size"))
	d.DiskSize = fmt.Sprintf("%dGi", flags.Int("harvester-disk-size"))
	d.DiskBus = flags.String("harvester-disk-bus")

	d.ImageName = flags.String("harvester-image-name")

	d.SSHUser = flags.String("harvester-ssh-user")
	d.SSHPort = flags.Int("harvester-ssh-port")

	d.KeyPairName = flags.String("harvester-key-pair-name")
	d.SSHPrivateKeyPath = flags.String("harvester-ssh-private-key-path")
	d.SSHPassword = flags.String("harvester-ssh-password")

	d.NetworkType = flags.String("harvester-network-type")

	d.NetworkName = flags.String("harvester-network-name")
	d.NetworkInterface = flags.String("harvester-network-interface")
	d.NetworkModel = flags.String("harvester-network-model")

	d.NetworkMask = flags.String("harvester-network-mask")
	d.NetworkGateway = flags.String("harvester-network-gateway")

	d.SetSwarmConfigFromFlags(flags)

	return d.checkConfig()
}

func (d *Driver) checkConfig() error {
	if d.Host == "" {
		return errors.New("must specify harvester host")
	}
	if d.Username == "" {
		return errors.New("must specify harvester username")
	}
	if d.Password == "" {
		return errors.New("must specify harvester password")
	}
	if d.ImageName == "" {
		return errors.New("must specify harvester image name")
	}
	if d.KeyPairName != "" && d.SSHPrivateKeyPath == "" {
		return errors.New("must specify the ssh private key path of the harvester key pair")
	}
	switch d.NetworkType {
	case networkTypePod:
	case networkTypeStatic, networkTypeDHCP:
		if d.NetworkName == "" {
			return errors.New("must specify harvester network name")
		}
		if d.NetworkInterface == "" {
			return errors.New("must specify harvester network interface")
		}
		if d.NetworkType == networkTypeStatic {
			if d.NetworkMask == "" {
				return errors.New("must specify harvester network mask")
			}
			if d.NetworkGateway == "" {
				return errors.New("must specify harvester network gateway")
			}
		}
	default:
		return fmt.Errorf("unknown network type %s", d.NetworkType)
	}
	return nil
}
