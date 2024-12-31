package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/undistro/cel-playground/eval"
	"io/ioutil"
	"testing"
)

type CelTestSuite struct {
	suite.Suite
	prog []byte
}

func (suite *CelTestSuite) SetupTest() {
	prog, err := ioutil.ReadFile("solution.txt")
	if err != nil {
		require.Failf(suite.T(), "Failed to read solution file: %s", err.Error())
	}
	suite.prog = prog
}

func (suite *CelTestSuite) TestPodSpecWithoutVolumes() {
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
	response, err := eval.CelEval(suite.prog, []byte(podSpecWithoutVolumes))
	require.Nil(suite.T(), err)
	var actual eval.EvalResponse
	err = json.Unmarshal([]byte(response), &actual)
	require.Nil(suite.T(), err)
	require.True(suite.T(), actual.Result.(bool))
}

func (suite *CelTestSuite) TestPodSpecWithNonGitRepoVolumesOnly() {
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
	response, err := eval.CelEval(suite.prog, []byte(podSpecWithNonGitRepoVolumesOnly))
	require.Nil(suite.T(), err)
	var actual eval.EvalResponse
	err = json.Unmarshal([]byte(response), &actual)
	require.Nil(suite.T(), err)
	require.True(suite.T(), actual.Result.(bool))
}

func (suite *CelTestSuite) TestPodSpecWithGitRepoVolumesOnly() {
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
	response, err := eval.CelEval(suite.prog, []byte(podSpecWithGitRepoVolumesOnly))
	require.Nil(suite.T(), err)
	var actual eval.EvalResponse
	err = json.Unmarshal([]byte(response), &actual)
	require.Nil(suite.T(), err)
	require.False(suite.T(), actual.Result.(bool))
}

func (suite *CelTestSuite) TestPodSpecWithGitRepoVolumes() {
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
	response, err := eval.CelEval(suite.prog, []byte(podSpecWithGitRepoVolumes))
	require.Nil(suite.T(), err)
	var actual eval.EvalResponse
	err = json.Unmarshal([]byte(response), &actual)
	require.Nil(suite.T(), err)
	require.False(suite.T(), actual.Result.(bool))
}

func TestForbidGitRepo(t *testing.T) {
	suite.Run(t, new(CelTestSuite))
}
