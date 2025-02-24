// A generated module for Dagger functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/dagger/internal/dagger"
)

type Dagger struct{}

// Release a new version of the project
func (m *Dagger) Release(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error){
	

	return dag.Node(dagger.NodeOpts{
		Version: "23.7.0"}).
		WithNpm().
		WithSource(source).
		Container().
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{"npm", "install","--save-dev","@semantic-release/git"}).
		WithExec([]string{"npm","install","--save-dev","@semantic-release/changelog"}).
		WithExec([]string{"npm","install","--save-dev","conventional-changelog-conventionalcommits"}).
		WithExec([]string{"npx","semantic-release"}).
	Stdout(ctx)
	
	// Commands().
	// Run([]string{"npx", "semantic-release", "--version"}).
  // Stdout(ctx)



	// version, err := dag.Nsv(source).Next(ctx, dagger.NsvNextOpts{Show: true}) 
	// if  err != nil {
	// 	return "", err
	// }
	//
	// fmt.Println("Version: ", version)
	//
	// return dag.Goreleaser(dagger.GoreleaserOpts{Source: source}).Test().Output(ctx)


	

}
