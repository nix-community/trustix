package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/libgit2/git2go/v30"
	"os"
	"time"
)

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
func (kv *GitKVStore) shardKey(key []byte) ([]byte, [][]byte) {
	var ret [][]byte
	hash := sha256.Sum256(key)

	for i := 0; i < kv.treeDepth; i++ {
		ret = append(ret, hash[kv.tokenLength*i:kv.tokenLength*i+kv.tokenLength])
	}
	ret = append(ret, hash[kv.tokenLength*kv.treeDepth:])

	return hash[:], ret
}

func (kv *GitKVStore) Set(key []byte, value []byte) error {
	hash, shardKey := kv.shardKey(key)
	hexPath := hex.EncodeToString(hash)

	fmt.Println(shardKey)

	odb, err := kv.repo.Odb()
	if err != nil {
		return err
	}
	defer odb.Free()

	builder, err := kv.repo.TreeBuilderFromTree(kv.tree)
	if err != nil {
		return err
	}
	defer builder.Free()

	blobId, err := odb.Write(value, git.ObjectBlob)
	if err != nil {
		return err
	}

	err = builder.Insert(hexPath, blobId, git.FilemodeBlob)
	if err != nil {
		return err
	}

	treeID, err := builder.Write()
	if err != nil {
		return err
	}

	tree, err := kv.repo.LookupTree(treeID)
	if err != nil {
		return err
	}

	kv.createCommit(fmt.Sprintf("Set %s", hexPath), tree)

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
	hash, _ := kv.shardKey(key)
	hexPath := hex.EncodeToString(hash)

	entry, err := kv.tree.EntryByPath(hexPath)
	if err != nil {
		return nil, err
	}

	blob, err := kv.repo.LookupBlob(entry.Id)
	if err != nil {
		return nil, err
	}
	defer blob.Free()

	return blob.Contents(), nil
}

func (kv *GitKVStore) Delete(key []byte) error {
	return nil
}
