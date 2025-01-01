package cel_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/undistro/cel-playground/eval"
)

// The user solution is available as userSolution
// var userSolution []byte

var _ = Describe("Forbid gitRepo volumes", func() {
	It("should work for a Pod spec without volumes", func() {
		podSpecWithoutVolumes := `
object:
  spec:
    containers:
      - args:
          - sleep
          - infinity
        image: busybox
        name: busybox
        resources: {}
    dnsPolicy: ClusterFirst
    restartPolicy: Always
`
		response, err := eval.CelEval(userSolution, []byte(podSpecWithoutVolumes))
		Expect(err).NotTo(HaveOccurred())
		var actual eval.EvalResponse
		Expect(json.Unmarshal([]byte(response), &actual)).To(Succeed())
		Expect(actual.Result.(bool)).To(BeTrue())
	})

	It("should work for a Pod spec with only non-gitRepo volumes", func() {
		podSpecWithNonGitRepoVolumesOnly := `
object:
  spec:
    containers:
      - args:
          - sleep
          - infinity
        image: busybox
        name: busybox
        resources: {}
        volumeMounts:
          - name: app-data
            mountPath: /data
    dnsPolicy: ClusterFirst
    restartPolicy: Always
    volumes:
      - name: app-data
        persistentVolumeClaim:
          claimName: busybox
`
		response, err := eval.CelEval(userSolution, []byte(podSpecWithNonGitRepoVolumesOnly))
		Expect(err).NotTo(HaveOccurred())
		var actual eval.EvalResponse
		Expect(json.Unmarshal([]byte(response), &actual)).To(Succeed())
		Expect(actual.Result.(bool)).To(BeTrue())
	})

	It("should work for a Pod spec with gitRepo volumes only", func() {
		podSpecWithGitRepoVolumesOnly := `
object:
  spec:
    containers:
      - args:
          - sleep
          - infinity
        image: busybox
        name: busybox
        resources: {}
        volumeMounts:
          - name: app-data
            mountPath: /data
    dnsPolicy: ClusterFirst
    restartPolicy: Always
    volumes:
      - name: app-data
        gitRepo:
          repository: "git@github.com:kubernetes/kubernetes.git"
          revision: "9fc9ddc7bceca86e805f674caff7d7acf31fad6c"
`
		response, err := eval.CelEval(userSolution, []byte(podSpecWithGitRepoVolumesOnly))
		Expect(err).NotTo(HaveOccurred())
		var actual eval.EvalResponse
		Expect(json.Unmarshal([]byte(response), &actual)).To(Succeed())
		Expect(actual.Result.(bool)).To(BeFalse())
	})

	It("should work for a Pod spec with gitRepo volumes and non-gitRepo volumes alike", func() {
		podSpecWithGitRepoVolumes := `
object:
  spec:
    containers:
      - args:
          - sleep
          - infinity
        image: busybox
        name: busybox
        resources: {}
        volumeMounts:
          - name: app-data
            mountPath: /data
          - name: app-config
            mountPath: /etc/my-app
    dnsPolicy: ClusterFirst
    restartPolicy: Always
    volumes:
      - name: app-data
        persistentVolumeClaim:
          claimName: busybox
      - name: app-config
        gitRepo:
          repository: "git@github.com:kubernetes/kubernetes.git"
          revision: "9fc9ddc7bceca86e805f674caff7d7acf31fad6c"
`
		response, err := eval.CelEval(userSolution, []byte(podSpecWithGitRepoVolumes))
		Expect(err).NotTo(HaveOccurred())
		var actual eval.EvalResponse
		Expect(json.Unmarshal([]byte(response), &actual)).To(Succeed())
		Expect(actual.Result.(bool)).To(BeFalse())
	})
})
