package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/libgit2/git2go/v30"
	"os"
	"time"
)

func dummy(_ interface{}) {

}

func insertNode(repo *git.Repository, treebuilder *git.TreeBuilder, path []string, content []byte) (*git.Oid, error) {

	odb, err := repo.Odb()
	if err != nil {
		return nil, err
	}
	defer odb.Free()

	treeOid, err := treebuilder.Write()
	if err != nil {
		panic(err)
	}
	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		panic(err)
	}
	// defer tree.Free()

	if len(path) == 1 {

		blobId, err := odb.Write(content, git.ObjectBlob)
		if err != nil {
			return nil, err
		}

		err = treebuilder.Insert(path[0], blobId, git.FilemodeBlob)
		if err != nil {
			return nil, err
		}

		return treebuilder.Write()
	}

	subtreeName := path[0]
	subPath := path[1:]

	var subTreebuilder *git.TreeBuilder
	entry := tree.EntryByName(subtreeName)
	if entry == nil {
		subTreebuilder, err = repo.TreeBuilder()
		if err != nil {
			panic(err)
		}
		// defer subTreebuilder.Free()
	} else {
		subTree, err := repo.LookupTree(entry.Id)
		if err != nil {
			panic(err)
		}

		subTreebuilder, err = repo.TreeBuilderFromTree(subTree)
		if err != nil {
			panic(err)
		}

	}

	subTreeOid, err := insertNode(repo, subTreebuilder, subPath, content)
	if err != nil {
		panic(err)
	}
	treebuilder.Insert(subtreeName, subTreeOid, git.FilemodeTree)
	return treebuilder.Write()

}

// GitKVStore - Use Git as a key/value store with automatic subtree sharding
type GitKVStore struct {
	repo *git.Repository

	// Sharding
	treeDepth   int
	tokenLength int

	// Commiter data
	name  string
	email string

	// Track Git state
	tree   *git.Tree
	commit *git.Commit // Previous commit
}

func newGitKVStore(repoPath string, name string, email string) (*GitKVStore, error) {

	// Always use bare repository (no worktree)
	bare := true

	// Hard code these for now, but should become configurable
	treeDepth := 5
	tokenLength := 2

	var repo *git.Repository
	var err error

	created := false

	if _, err = os.Stat(repoPath); os.IsNotExist(err) {
		created = true
		// Repo doesn't exist, create it
		repo, err = git.InitRepository(repoPath, bare)
		if err != nil {
			return nil, err
		}
	} else {
		repo, err = git.OpenRepository(repoPath)
		if err != nil {
			return nil, err
		}
	}

	instance := &GitKVStore{
		repo:        repo,
		treeDepth:   treeDepth,
		tokenLength: tokenLength,
		name:        name,
		email:       email,
	}

	if created {
		builder, err := repo.TreeBuilder()
		if err != nil {
			return nil, err
		}
		defer builder.Free()

		treeOid, err := builder.Write()
		if err != nil {
			return nil, err
		}

		tree, err := repo.LookupTree(treeOid)
		if err != nil {
			return nil, err
		}
		defer tree.Free()

		sig := instance.createSig()
		message := "Init"

		commitOid, err := repo.CreateCommit("HEAD", sig, sig, message, tree)
		if err != nil {
			return nil, err
		}

		commit, err := repo.LookupCommit(commitOid)
		if err != nil {
			return nil, err
		}
		defer commit.Free()
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	defer head.Free()

	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	instance.tree = tree
	instance.commit = commit

	return instance, nil

}

func (kv *GitKVStore) createSig() *git.Signature {
	return &git.Signature{
		Name:  kv.name,
		Email: kv.email,
		When:  time.Now(),
	}
}

// shardKey - Hash a key and return the path to the object in the Git tree
func (kv *GitKVStore) shardKey(key []byte) []string {
	var ret []string
	bH := sha256.Sum256(key)
	hash := hex.EncodeToString(bH[:])

	for i := 0; i < kv.treeDepth; i++ {
		ret = append(ret, hash[kv.tokenLength*i:kv.tokenLength*i+kv.tokenLength])
	}
	ret = append(ret, hash[kv.tokenLength*kv.treeDepth:])

	return ret
}

func (kv *GitKVStore) Set(key []byte, value []byte) error {
	shardKey := kv.shardKey(key)

	builder, err := kv.repo.TreeBuilderFromTree(kv.tree)
	if err != nil {
		return err
	}
	// defer builder.Free()

	treeOid, err := insertNode(kv.repo, builder, shardKey, value)
	if err != nil {
		return err
	}

	tree, err := kv.repo.LookupTree(treeOid)
	if err != nil {
		return err
	}

	err = kv.createCommit("something", tree)
	if err != nil {
		return err
	}

	return nil
}

func (kv *GitKVStore) createCommit(message string, tree *git.Tree) error {

	sig := kv.createSig()

	commitOid, err := kv.repo.CreateCommit("HEAD", sig, sig, message, tree, kv.commit)
	if err != nil {
		return err
	}

	commit, err := kv.repo.LookupCommit(commitOid)
	if err != nil {
		return err
	}

	oldCommit := commit
	oldTree := kv.tree

	kv.tree = tree
	kv.commit = commit

	oldTree.Free()
	oldCommit.Free()

	return nil
}

func (kv *GitKVStore) Get(key []byte) ([]byte, error) {
	shardKey := kv.shardKey(key)

	tree := kv.tree

	for i, p := range shardKey {
		var err error
		entry := tree.EntryByName(p)
		if entry == nil {
			return nil, fmt.Errorf("Path component %s not found in tree", p)
		}

		if i+1 == len(shardKey) {
			blob, err := kv.repo.LookupBlob(entry.Id)
			if err != nil {
				return nil, err
			}

			return blob.Contents(), nil
		}

		tree, err = kv.repo.LookupTree(entry.Id)
		if err != nil {
			panic(err)
		}

	}

	return nil, fmt.Errorf("NOPO")
}

func (kv *GitKVStore) Delete(key []byte) error {
	return nil
}
