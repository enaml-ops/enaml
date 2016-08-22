package generators_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	. "github.com/enaml-ops/enaml/generators"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/lo"
)

var _ = Describe("generate", func() {
	testJobGeneration("./fixtures/gorouter-job-duplicates.MF", "gorouter", 11, []string{"clients.go", "gorouter.go", "gorouterjob.go", "logrotate.go", "metron.go", "nats.go", "router.go", "routingapi.go", "ssl.go", "status.go", "uaa.go"})
	testJobGeneration("./fixtures/atc-job.MF", "atc", 8, []string{"atcjob.go", "container.go", "githubauth.go", "postgresql.go", "retention.go", "riemann.go", "role.go", "yeller.go"})
	testJobGeneration("./fixtures/haproxy-job.MF", "haproxy", 6, []string{"appssh.go", "cc.go", "haproxy.go", "haproxyjob.go", "router.go", "servers.go"})
	testJobGeneration("./fixtures/gorouter-job.MF", "gorouter", 11, []string{"clients.go", "gorouter.go", "gorouterjob.go", "logrotate.go", "metron.go", "nats.go", "router.go", "routingapi.go", "ssl.go", "status.go", "uaa.go"})
	testJobGeneration("./fixtures/uaa-job.MF", "uaa", 43, []string{"admin.go",
		"analytics.go",
		"authentication.go",
		"authenticationpolicy.go",
		"branding.go",
		"claims.go",
		"client.go",
		"database.go",
		"env.go",
		"global.go",
		"internal.go",
		"jwt.go",
		"jwtpolicy.go",
		"ldapgroups.go",
		"ldapssl.go",
		"links.go",
		"login.go",
		"loginldap.go",
		"logout.go",
		"notifications.go",
		"oauth.go",
		"parameter.go",
		"passwordpolicy.go",
		"prompt.go",
		"promptpassword.go",
		"proxy.go",
		"redirect.go",
		"saml.go",
		"scim.go",
		"scimuser.go",
		"servlet.go",
		"smtp.go",
		"socket.go",
		"uaa.go",
		"uaadb.go",
		"uaajob.go",
		"uaaldap.go",
		"uaalogin.go",
		"uaapassword.go",
		"uaassl.go",
		"uaauser.go",
		"username.go",
		"zones.go"})
})

func testJobGeneration(controlRelease, packageName string, numberOfFiles int, expectedFileNames []string) {
	Describe(fmt.Sprintf("given a Job definition from %s", controlRelease), func() {
		var pwd, _ = os.Getwd()
		var tmpDirName = "_testGen"
		var tmpPath = path.Join(pwd, tmpDirName)
		var controlOutputDir string
		BeforeEach(func() {
			os.MkdirAll(tmpPath, 0700)
			controlOutputDir, _ = ioutil.TempDir(tmpPath, "releasejobs")
		})

		AfterEach(func() {
			lo.G.Debug("removing: ", controlOutputDir)
			os.RemoveAll(tmpPath)
		})

		Context("when called on a release", func() {
			BeforeEach(func() {
				var bytes []byte
				var err error

				if bytes, err = ioutil.ReadFile(controlRelease); err != nil {
					fmt.Println(err)
					panic(err)
				}
				Generate(packageName, bytes, controlOutputDir)
			})
			It(fmt.Sprintf("then it should create %v", numberOfFiles), func() {
				files, _ := ioutil.ReadDir(controlOutputDir)
				Ω(len(files)).Should(Equal(numberOfFiles))
			})
			It("then it should create files with correct names", func() {
				fileNames := make([]string, 0)
				files, _ := ioutil.ReadDir(controlOutputDir)

				for _, file := range files {
					fileNames = append(fileNames, file.Name())
				}
				Ω(fileNames).Should(ConsistOf(expectedFileNames))
			})

			It("then it should create a struct with the correct elements and format for each nested property", func() {
				var errBuffer bytes.Buffer
				lo.G.Debug(controlOutputDir)
				lo.G.Debug("package", packageName)
				cmdArgs := "./" + path.Join(tmpDirName, path.Base(controlOutputDir))
				fmt.Println(cmdArgs)
				cmd := exec.Command("go", "build", cmdArgs)
				cmd.Stderr = &errBuffer
				out, err := cmd.Output()
				lo.G.Debug("out: ", out)
				Ω(errBuffer.String()).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
}
