package init_new

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/plugins/project/validators"
	"github.com/Red-Sock/rscli/tests/mocks"
)

const (
	defaultGitPath = "github.com/RedSock"
	projName       = "test_proj"
)

func Test_collectName(t *testing.T) {

	type testCase struct {
		userInput []string
		args      []string

		io            *mocks.IOMock
		printMockCall func(in string)
		expectedErr   error
		expectedResp  string
	}

	type testCaseConstructor struct {
		new func() testCase
	}

	testCases := map[string]testCaseConstructor{
		"OK_SHORT_NAME_ARG": {
			new: func() (tc testCase) {

				tc.args = []string{projName}
				tc.io = mocks.NewIOMock(t)

				tc.expectedResp = defaultGitPath + "/" + projName

				tc.io.PrintlnColoredMock.Expect(colors.ColorCyan, fmt.Sprintf(ackProjectNameMessagePattern, tc.expectedResp))

				return
			},
		},
		"OK_SHORT_NAME_USER_INPUT": {
			new: func() (tc testCase) {
				const projName = "test_proj"

				tc.io = mocks.NewIOMock(t)

				tc.io.PrintMock.Expect(fmt.Sprintf(askUserForNameMessagePattern, defaultGitPath))
				tc.io.GetInputMock.Expect().Return(projName, nil)

				tc.expectedResp = defaultGitPath + "/" + projName

				tc.io.PrintlnColoredMock.Expect(colors.ColorCyan, fmt.Sprintf(ackProjectNameMessagePattern, tc.expectedResp))

				return
			},
		},

		"OK_LONG_NAME_ARG": {
			new: func() (tc testCase) {

				tc.expectedResp = defaultGitPath + "/" + projName

				tc.args = []string{tc.expectedResp}
				tc.io = mocks.NewIOMock(t)

				tc.io.PrintlnColoredMock.Expect(colors.ColorCyan, fmt.Sprintf(ackProjectNameMessagePattern, tc.expectedResp))

				return
			},
		},

		"OK_LONG_NAME_WITH_SEC_PROTOC_ARG": {
			new: func() (tc testCase) {

				tc.expectedResp = defaultGitPath + "/" + projName

				tc.args = []string{"https://" + tc.expectedResp}
				tc.io = mocks.NewIOMock(t)

				tc.io.PrintlnColoredMock.Expect(colors.ColorCyan, fmt.Sprintf(ackProjectNameMessagePattern, tc.expectedResp))

				return
			},
		},
		"OK_LONG_NAME_WITH_PLAIN_PROTOC_ARG": {
			new: func() (tc testCase) {

				tc.expectedResp = defaultGitPath + "/" + projName

				tc.args = []string{"http://" + tc.expectedResp}
				tc.io = mocks.NewIOMock(t)

				tc.io.PrintlnColoredMock.Expect(colors.ColorCyan, fmt.Sprintf(ackProjectNameMessagePattern, tc.expectedResp))

				return
			},
		},

		"ERR_BOTH_EMPTY": {
			new: func() (tc testCase) {
				tc.io = mocks.NewIOMock(t)
				tc.io.PrintMock.Expect(fmt.Sprintf(askUserForNameMessagePattern, defaultGitPath))
				tc.io.GetInputMock.Expect().Return("", nil)
				tc.expectedErr = emptyNameErr
				return
			},
		},
		"ERR_GETTING_USER_RESPONSE": {
			new: func() (tc testCase) {
				tc.expectedErr = rerrors.New("some err")

				tc.io = mocks.NewIOMock(t)
				tc.io.PrintMock.Expect(fmt.Sprintf(askUserForNameMessagePattern, defaultGitPath))
				tc.io.GetInputMock.Expect().Return("", tc.expectedErr)
				return
			},
		},
		"ERR_FAILED_NAME_VALIDATION": {
			new: func() (tc testCase) {
				tc.expectedErr = validators.ErrInvalidNameErr
				tc.args = []string{"%"}
				tc.io = mocks.NewIOMock(t)
				tc.io.GetInputMock.Expect().Return("", tc.expectedErr)
				return
			},
		},
	}

	for name, constructor := range testCases {
		name, constructor := name, constructor
		t.Run(name, func(t *testing.T) {
			tc := constructor.new()
			nc := newNameCollector(tc.io, defaultGitPath)

			resp, err := nc.collect(tc.args)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, resp, tc.expectedResp)
		})
	}

}
