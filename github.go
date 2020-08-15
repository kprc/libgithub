package libgithub

import (
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type RepoSetting struct {
	Owner string
	Repository string
	Path string
	Token string
	Name string
	Email string
}

type GithubClient struct {
	repo RepoSetting
	client *github.Client
}


//http://github.com/kprc/libgithub/github.go
//token access token from your setting
//owner  kprc
//repo  libgithub
//filePath github.go
func NewGithubClient(token, owner, repo, filePath, name, email string) *GithubClient  {
	gc := &GithubClient{}
	gc.repo.Path = filePath
	gc.repo.Owner = owner
	gc.repo.Token = token
	gc.repo.Repository = repo
	gc.repo.Name = name
	gc.repo.Email = email

	return gc
}

func (gc *GithubClient)CreateFile(commitMsg ,fileContent string) error  {
	if gc.client == nil{
		ts:=oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gc.repo.Token})
		tc:=oauth2.NewClient(context.Background(),ts)
		gc.client = github.NewClient(tc)
	}

	rcfo:=&github.RepositoryContentFileOptions{}
	rcfo.Message = &commitMsg
	rcfo.Content = []byte(fileContent)
	rcfo.Committer = &github.CommitAuthor{
		Name: &gc.repo.Name,
		Email: &gc.repo.Email,
	}

	repoResp,resp,err:=gc.client.Repositories.CreateFile(context.Background(),gc.repo.Owner,gc.repo.Repository,gc.repo.Path,rcfo)
	if err!=nil{
		return err
	}
	if repoResp != nil{
		fmt.Println(*repoResp)
	}
	if resp != nil{
		fmt.Println(*resp)
	}

	return nil
}

func (gc *GithubClient)UpdateFile(commitMsg, fileContent string) error  {
	if gc.client == nil{
		ts:=oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gc.repo.Token})
		tc:=oauth2.NewClient(context.Background(),ts)
		gc.client = github.NewClient(tc)
	}

	_,hash,err:=gc.GetContent()
	if err!=nil{
		return err
	}

	rcfo:=&github.RepositoryContentFileOptions{}
	rcfo.Message = &commitMsg
	rcfo.Content = []byte(fileContent)
	rcfo.Committer = &github.CommitAuthor{
		Name: &gc.repo.Name,
		Email: &gc.repo.Email,
	}
	rcfo.SHA = &hash

	var (
		repoResp *github.RepositoryContentResponse
		resp *github.Response
	)
	repoResp,resp,err=gc.client.Repositories.UpdateFile(context.Background(),gc.repo.Owner,gc.repo.Repository,gc.repo.Path,rcfo)
	if err!=nil{
		return err
	}
	if repoResp != nil{
		fmt.Println(*repoResp)
	}
	if resp != nil{
		fmt.Println(*resp)
	}

	return nil
}

func (gc *GithubClient)GetContent() (content ,hash string, err error)  {
	if gc.client == nil{
		ts:=oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gc.repo.Token})
		tc:=oauth2.NewClient(context.Background(),ts)
		gc.client = github.NewClient(tc)
	}
	fc,d,resp,err:=gc.client.Repositories.GetContents(context.Background(),gc.repo.Owner,gc.repo.Repository,gc.repo.Path,nil)
	if err!=nil{
		return "","",err
	}
	if fc!=nil{
		fmt.Println(*fc)
	}
	if d!=nil{
		fmt.Println(d)
	}
	if resp != nil{
		fmt.Println(*resp)
	}
	var plaintxt []byte
	plaintxt,err = base64.StdEncoding.DecodeString(*fc.Content)

	hash = fc.GetSHA()

	return string(plaintxt),hash,nil
}




