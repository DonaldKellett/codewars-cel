package main

import (
	"encoding/json"
	"fmt"
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
		panic(fmt.Sprintf("Failed to read solution file: %v", err))
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

func (suite *CelTestSuite) TestPodSpecWithNonGitRepoVolumes() {
	podSpecWithNonGitRepoVolumes := `
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
	response, err := eval.CelEval(suite.prog, []byte(podSpecWithNonGitRepoVolumes))
	require.Nil(suite.T(), err)
	var actual eval.EvalResponse
	err = json.Unmarshal([]byte(response), &actual)
	require.Nil(suite.T(), err)
	require.True(suite.T(), actual.Result.(bool))
}

func TestForbidGitRepo(t *testing.T) {
	suite.Run(t, new(CelTestSuite))
}
