package e2e_test

import (
	"io/ioutil"
	"os/exec"

	"github.com/hyperledger-twgc/tape/e2e"
	"github.com/hyperledger-twgc/tape/e2e/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Mock test for good path", func() {

	Context("E2E with multi mocked Fabric", func() {
		When("traffic and observer mode", func() {

			It("should work properly", func() {
				server, err := mock.NewServer(2, nil)
				server.Start()
				defer server.Stop()

				config, err := ioutil.TempFile("", "endorsement-only-config-*.yaml")
				paddrs, oaddr := server.Addresses()
				configValue := e2e.Values{
					PrivSk:          mtlsKeyFile.Name(),
					SignCert:        mtlsCertFile.Name(),
					Mtls:            false,
					PeersAddrs:      paddrs,
					OrdererAddr:     oaddr,
					CommitThreshold: 1,
				}
				e2e.GenerateConfigFile(config.Name(), configValue)

				cmd0 := exec.Command(tapeBin, "traffic", "-c", config.Name(), "--rate=10")

				//cmd1 := exec.Command(tapeBin, "observer", "-c", config.Name())
				tapeSession, err = gexec.Start(cmd0, nil, nil)
				Expect(err).NotTo(HaveOccurred())
				//_, err = gexec.Start(cmd0, nil, nil)
				//Expect(err).NotTo(HaveOccurred())
				//Eventually(tapeSession.Out).Should(Say("Time.*Tx.*"))
			})

			It("should work properly", func() {
				server, err := mock.NewServer(2, nil)
				server.Start()
				defer server.Stop()

				config, err := ioutil.TempFile("", "endorsement-only-config-*.yaml")
				paddrs, oaddr := server.Addresses()
				configValue := e2e.Values{
					PrivSk:          mtlsKeyFile.Name(),
					SignCert:        mtlsCertFile.Name(),
					Mtls:            false,
					PeersAddrs:      paddrs,
					OrdererAddr:     oaddr,
					CommitThreshold: 1,
					PolicyFile:      PolicyFile.Name(),
				}
				e2e.GenerateConfigFile(config.Name(), configValue)

				cmd0 := exec.Command(tapeBin, "traffic", "-c", config.Name())

				cmd1 := exec.Command(tapeBin, "observer", "-c", config.Name())
				tapeSession, err = gexec.Start(cmd1, nil, nil)
				Expect(err).NotTo(HaveOccurred())
				_, err = gexec.Start(cmd0, nil, nil)
				Expect(err).NotTo(HaveOccurred())
				Eventually(tapeSession.Out).Should(Say("Time.*Tx.*"))
			})
		})
	})
})
