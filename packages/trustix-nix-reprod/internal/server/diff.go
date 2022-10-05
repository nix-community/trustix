// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	connect "github.com/bufbuild/connect-go"
	"github.com/nix-community/go-nix/pkg/nar"
	"github.com/nix-community/go-nix/pkg/nixpath"
	"github.com/nix-community/trustix/packages/go-lib/executor"
	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/future"
	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/refcount"
	pb "github.com/nix-community/trustix/packages/trustix-nix-reprod/reprod-api"
	"github.com/nix-community/trustix/packages/trustix-nix/schema"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix/client"
)

type narDownload struct {
	narinfo *schema.NarInfo
	narPath string // Path to NAR on the local file system
}

var epoch = time.Unix(1, 0)

func downloadNAR(ctx context.Context, client *client.Client, outputHash string) (*refcount.RefCountedValue[*narDownload], error) {
	narinfo := &schema.NarInfo{}
	{
		digest, err := base64.URLEncoding.DecodeString(outputHash)
		if err != nil {
			return nil, err
		}

		resp, err := client.NodeAPI.GetValue(ctx, &api.ValueRequest{
			Digest: digest,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting log entries: %w", err)
		}

		err = json.Unmarshal(resp.Value, narinfo)
		if err != nil {
			return nil, fmt.Errorf("error decoding narinfo: %w", err)
		}

		if !strings.HasPrefix(narinfo.StorePath, nixpath.StoreDir) {
			return nil, fmt.Errorf("path '%s' not starting with store prefix '%s'", narinfo.StorePath, nixpath.StoreDir)
		}
	}

	storePrefixHash := narinfo.StorePath[len(nixpath.StoreDir)+1 : len(nixpath.StoreDir)+32+1]

	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		return nil, err
	}

	u.Path = path.Join("nar", storePrefixHash, url.QueryEscape(narinfo.NarHash))

	out, err := os.CreateTemp("", "trustix-nix-reprod-nar")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	return refcount.NewRefCountedValue(&narDownload{
		narinfo: narinfo,
		narPath: out.Name(),
	}, func() error {
		// Delete temporary nar file when all references have been dropped
		return os.Remove(out.Name())
	}), nil
}

func unpackNar(narpath string, unpackDir string) error {
	f, err := os.Open(narpath)
	if err != nil {
		return fmt.Errorf("error opening '%s': %w", narpath, err)
	}
	defer f.Close()

	r, err := nar.NewReader(f)
	if err != nil {
		return fmt.Errorf("error creating reader: %w", err)
	}
	defer r.Close()

	// Make dirs read-only when done unpacking
	dirs := []string{}

	// Reset file times when done on packing
	files := []string{}

	handleHeader := func(header *nar.Header) error {
		storePath := path.Join(unpackDir, header.Path)

		files = append(files, storePath)

		switch header.Type {
		case nar.TypeRegular:
			mode := os.FileMode(0444)
			if header.Executable {
				mode = mode | 0111
			}

			f, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
			if err != nil {
				return fmt.Errorf("error creating file '%s': %w", storePath, err)
			}
			defer f.Close()

			_, err = io.Copy(f, r)
			if err != nil {
				return fmt.Errorf("error copying file data to '%s': %w", storePath, err)
			}

		case nar.TypeDirectory:
			err = os.Mkdir(storePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory '%s': %w", storePath, err)
			}

			dirs = append(dirs, storePath)

		case nar.TypeSymlink:
			err = os.Symlink(header.LinkTarget, storePath)
			if err != nil {
				return fmt.Errorf("error creating directory '%s': %w", storePath, err)
			}

		default:
			return fmt.Errorf("unhandled nar header type: %v", header.Type)
		}

		return nil
	}

	for {
		header, err := r.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("error getting NAR header: %w", err)
		}

		err = handleHeader(header)
		if err != nil {
			return err
		}
	}

	for _, dir := range dirs {
		err = os.Chmod(dir, 0555)
		if err != nil {
			return fmt.Errorf("error setting directory '%s' to read only: '%w'", dir, err)
		}
	}

	for _, file := range files {
		err := os.Chtimes(file, epoch, epoch)
		if err != nil {
			return fmt.Errorf("error setting file times for '%s': %w", file, err)
		}
	}

	return nil
}

func downloadAndUnpackStorePath(downloadExecutor *future.KeyedFutures[*refcount.RefCountedValue[*narDownload]], client *client.Client, outputHash string, tmpDir string, unpackDirSuffix string) (string, error) {

	fut := downloadExecutor.Run(outputHash, func() (*refcount.RefCountedValue[*narDownload], error) {
		return downloadNAR(context.Background(), client, outputHash)
	})

	ref, err := fut.Result()
	if err != nil {
		return "", fmt.Errorf("error downloading NAR with output hash '%s': %w", outputHash, err)
	}

	ref.Incr() // Nar download reference
	defer func() {
		err = ref.Decr()
		if err != nil {
			panic(fmt.Errorf("error cleaning up reference counted value: %w", err))
		}
	}()

	if !strings.HasPrefix(ref.Value.narinfo.StorePath, nixpath.StoreDir) {
		return "", fmt.Errorf("store path '%s' outside of store dir prefix '%s'", ref.Value.narinfo.StorePath, nixpath.StoreDir)
	}

	storeDir := path.Join(tmpDir, ref.Value.narinfo.StorePath[1:len(ref.Value.narinfo.StorePath)])
	unpackDir := path.Join(storeDir, unpackDirSuffix)

	err = os.MkdirAll(storeDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating '%s': %w", storeDir, err)
	}

	// Unpack the guy into temp workspace
	err = unpackNar(ref.Value.narPath, unpackDir)
	if err != nil {
		return "", fmt.Errorf("error unpacking nar: %w", err)
	}

	return unpackDir, nil
}

func diff(downloadExecutor *future.KeyedFutures[*refcount.RefCountedValue[*narDownload]], db *sql.DB, client *client.Client, outputHash1 string, outputHash2 string) (*pb.DiffResponse, error) {
	outputHashes := []string{outputHash1, outputHash2}
	sort.Strings(outputHashes) // canonicalise output (same output no matter argument ordering)

	tmpDir, err := os.MkdirTemp("", "trustix-nix-reprod-diff")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary store dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	comparePaths := make([]string, len(outputHashes))

	// Download/unpack NARs
	{
		e := executor.NewParallellExecutor()

		for i, outputHash := range outputHashes {
			i := i
			outputHash := outputHash

			var unpackDirSuffix string
			switch i {
			case 0:
				unpackDirSuffix = "A"
			case 1:
				unpackDirSuffix = "B"
			default:
				panic(fmt.Errorf("unexpected index: %d", i))
			}

			e.Add(func() error {
				unpackedDir, err := downloadAndUnpackStorePath(downloadExecutor, client, outputHash, tmpDir, unpackDirSuffix)
				if err != nil {
					return err
				}

				comparePaths[i] = unpackedDir

				return nil
			})
		}

		err = e.Wait()
		if err != nil {
			return nil, fmt.Errorf("error downloading/unpacking NARs: %w", err)
		}
	}

	resp := &pb.DiffResponse{}

	// Run diffoscope
	{
		dirARel := strings.TrimPrefix(comparePaths[0], tmpDir+"/")
		dirBRel := strings.TrimPrefix(comparePaths[1], tmpDir+"/")

		cmd := exec.Command("diffoscope", "--html", "-", dirARel, dirBRel)
		cmd.Dir = tmpDir

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Diffoscope returns non-zero on paths that have a diff
		// Instead use stderr as a heurestic if the call went well or not
		cmd.Run() // Ignore error

		if stderr.Len() > 0 {
			if err != nil {
				return nil, fmt.Errorf("error executing diffoscope: %s", stderr.String())
			}
		}

		resp.HTMLDiff = stdout.String()
	}

	return resp, nil
}

func (s *APIServer) Diff(ctx context.Context, req *connect.Request[pb.DiffRequest]) (*connect.Response[pb.DiffResponse], error) {
	msg := req.Msg

	outputHash1 := msg.OutputHash1
	outputHash2 := msg.OutputHash2

	requestKey := outputHash1 + "." + outputHash2
	if outputHash2 > outputHash1 {
		requestKey = outputHash2 + "." + outputHash1
	}

	resp, err := s.diffExecutor.Run(requestKey, func() (*pb.DiffResponse, error) {
		return diff(s.downloadExecutor, s.db, s.client, outputHash1, outputHash2)
	}).Result()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}
