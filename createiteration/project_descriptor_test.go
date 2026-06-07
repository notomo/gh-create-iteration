package createiteration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjectDescriptor(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			name       string
			projectUrl string
			want       ProjectDescriptor
		}{
			{
				name:       "user project url",
				projectUrl: "https://github.com/users/notomo/projects/1",
				want: ProjectDescriptor{
					OwnerIsOrganization: false,
					Owner:               "notomo",
					Number:              1,
				},
			},
			{
				name:       "organization project url",
				projectUrl: "https://github.com/orgs/notomoorg/projects/2",
				want: ProjectDescriptor{
					OwnerIsOrganization: true,
					Owner:               "notomoorg",
					Number:              2,
				},
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				got, err := GetProjectDescriptor(c.projectUrl)
				require.NoError(t, err)
				assert.Equal(t, c.want, *got)
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := GetProjectDescriptor("https://github.com/notomo/example/issues/1")
		require.Error(t, err)
	})
}
