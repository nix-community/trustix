package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/libgit2/git2go/v30"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/storage/errors"
	"os"
	"time"
)

type gitTxn struct {
	tree *git.Tree
	repo *git.Repository
}

func (t *gitTxn) shardKey(key []byte) []string {
	// HACK: Special handling of HEAD (cheaper lookup)
	if bytes.Compare(key, []byte("HEAD")) == 0 {
		return []string{"HEAD"}
	}

	treeDepth := 5
	tokenLength := 2

	var ret []string
	bH := sha256.Sum256(key)
	hash := hex.EncodeToString(bH[:])

	for i := 0; i < treeDepth; i++ {
		ret = append(ret, hash[tokenLength*i:tokenLength*i+tokenLength])
	}
	ret = append(ret, hash[tokenLength*treeDepth:])

	return ret
}

func (t *gitTxn) Get(key []byte) ([]byte, error) {
	path := t.shardKey(key)

	tree := t.tree
	for i, p := range path {
		var err error
		entry := tree.EntryByName(p)
		if entry == nil {
			return nil, errors.ObjectNotFoundError
		}

		if i+1 == len(path) {
			blob, err := t.repo.LookupBlob(entry.Id)
			if err != nil {
				return nil, err
			}

			return blob.Contents(), nil
		}

		tree, err = t.repo.LookupTree(entry.Id)
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.ObjectNotFoundError
}

func (t *gitTxn) Set(key []byte, value []byte) error {
	path := t.shardKey(key)

	builder, err := t.repo.TreeBuilderFromTree(t.tree)
	if err != nil {
		return err
	}
	defer builder.Free()

	treeOid, err := insertGitNode(t.repo, builder, path, value)
	if err != nil {
		return err
	}

	tree, err := t.repo.LookupTree(treeOid)
	if err != nil {
		return err
	}

	t.tree = tree

	return nil
}

func (t *gitTxn) commit() error {
	return nil
}

func insertGitNode(repo *git.Repository, treebuilder *git.TreeBuilder, path []string, content []byte) (*git.Oid, error) {
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
	defer tree.Free()

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
		defer subTreebuilder.Free()
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

	subTreeOid, err := insertGitNode(repo, subTreebuilder, subPath, content)
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

func GitStorageFromConfig(conf *config.GitStorageConfig) (*GitKVStore, error) {

	// Always use bare repository (no worktree)
	bare := true

	// Hard code these for now, but should become configurable
	treeDepth := 5
	tokenLength := 2

	var repo *git.Repository
	var err error

	created := false

	if _, err = os.Stat(conf.Path); os.IsNotExist(err) {
		created = true
		// Repo doesn't exist, create it
		repo, err = git.InitRepository(conf.Path, bare)
		if err != nil {
			return nil, err
		}
	} else {
		repo, err = git.OpenRepository(conf.Path)
		if err != nil {
			return nil, err
		}
	}

	instance := &GitKVStore{
		repo:        repo,
		treeDepth:   treeDepth,
		tokenLength: tokenLength,
		name:        conf.Commiter,
		email:       conf.Email,
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

func (kv *GitKVStore) Close() {
}

func (kv *GitKVStore) runTX(readWrite bool, fn func(Transaction) error) error {
	t := &gitTxn{
		tree: kv.tree,
		repo: kv.repo,
	}

	err := fn(t)
	if err != nil {
		return err
	} else {
		if readWrite {
			return kv.createCommit("commit", t.tree)
			// return t.commit()
		}
	}

	return err
}

func (kv *GitKVStore) View(fn func(Transaction) error) error {
	return kv.runTX(false, fn)
}

func (kv *GitKVStore) Update(fn func(Transaction) error) error {
	return kv.runTX(true, fn)
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
	kv.commit = commit
	kv.tree = tree
	oldCommit.Free()

	return nil
}
