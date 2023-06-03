package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCMDValidate(t *testing.T) {
	testCases := []struct {
		name      string
		cmd       RunCMD
		whiteList []string
		err       string
	}{
		{
			name: "happy-path",
			cmd: RunCMD{
				DockerImage:   "happy-path",
				ContainerName: "happy-path",
			},
			whiteList: []string{"happy-path"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cmd.Validate(tc.whiteList...)
			if tc.err == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tc.err)
		})
	}
}
