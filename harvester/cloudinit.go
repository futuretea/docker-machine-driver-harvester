package harvester

import "fmt"

const (
	userDataHeader            = `#cloud-config`
	userDataAddQemuGuestAgent = `
package_update: true
packages:
- qemu-guest-agent
runcmd:
- [systemctl, enable, --now, qemu-guest-agent]`
	userDataPasswordTemplate = `
user: %s
password: %s
chpasswd: { expire: False }
ssh_pwauth: True`

	userDataSSHKeyTemplate = `
ssh_authorized_keys:
- >-
  %s`
	userDataAddDockerGroupSSHKeyTemplate = `
groups:
- docker
users:
- name: %s
  sudo: ALL=(ALL) NOPASSWD:ALL
  groups: sudo, docker
  shell: /bin/bash
  ssh_authorized_keys:
  - >-
    %s`
	networkDataTemplate = `
network:
  version: 1
  config:
  - type: physical
    name: %s`
	networkDataDHCPTemplate = `
    subnets:
    - type: dhcp`
	networkDataStaticTemplate = `
    subnets:
    - type: static
      address: %s
      netmask: %s
      gateway: %s`
)

func (d *Driver) createCloudInit() (userData string, networkData string) {
	// userData
	userData = userDataHeader
	if d.NetworkInterface != "" {
		// need qemu guest agent to get ip
		userData += userDataAddQemuGuestAgent
	}
	if d.SSHPassword != "" {
		userData += fmt.Sprintf(userDataPasswordTemplate, d.SSHUser, d.SSHPassword)
	}
	if d.SSHPublicKey != "" {
		if d.AddUserToDockerGroup {
			userData += fmt.Sprintf(userDataAddDockerGroupSSHKeyTemplate, d.SSHUser, d.SSHPublicKey)
		} else {
			userData += fmt.Sprintf(userDataSSHKeyTemplate, d.SSHPublicKey)
		}
	}
	// networkData
	if d.NetworkInterface != "" {
		networkData = fmt.Sprintf(networkDataTemplate, d.NetworkInterface)
		if d.IPAddress != "" && d.NetworkGateway != "" && d.NetworkMask != "" {
			networkData += fmt.Sprintf(networkDataStaticTemplate, d.IPAddress, d.NetworkMask, d.NetworkGateway)
		} else {
			networkData += networkDataDHCPTemplate
		}
	}
	return
}
