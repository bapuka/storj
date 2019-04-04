// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"storj.io/storj/internal/fpath"
	libuplink "storj.io/storj/lib/uplink"
	"storj.io/storj/pkg/process"
	"storj.io/storj/pkg/storj"
)

var (
	recursiveFlag *bool
)

func init() {
	lsCmd := addCmd(&cobra.Command{
		Use:   "ls",
		Short: "List objects and prefixes or all buckets",
		RunE:  list,
	}, RootCmd)
	recursiveFlag = lsCmd.Flags().Bool("recursive", false, "if true, list recursively")
}

func list(cmd *cobra.Command, args []string) error {
	ctx := process.Ctx(cmd)

	project, err := cfg.GetProject(ctx)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		src, err := fpath.New(args[0])
		if err != nil {
			return err
		}

		if src.IsLocal() {
			return fmt.Errorf("No bucket specified, use format sj://bucket/")
		}

		if err != nil {
			return err
		}

		enc := cfg.GetEncryptionScheme()

		bucket, err := project.OpenBucket(ctx, src.String(), &libuplink.EncryptionAccess{
			Key: enc.Cipher,
		}, 0)

		if err != nil {
			return err
		}
		err = listFiles(ctx, bucket, src, false)

		return convertError(err, src)
	}

	startAfter := ""
	noBuckets := true

	for {
		list, err := project.ListBuckets(ctx, &storj.BucketListOptions{Direction: storj.After, Cursor: startAfter})
		if err != nil {
			return err
		}
		if len(list.Items) > 0 {
			noBuckets = false
			for _, bucket := range list.Items {
				fmt.Println("BKT", formatTime(bucket.Created), bucket.Name)
				if *recursiveFlag {
					prefix, err := fpath.New(fmt.Sprintf("sj://%s/", bucket.Name))
					if err != nil {
						return err
					}
					err = listFiles(ctx, project, prefix, true)
					if err != nil {
						return err
					}
				}
			}
		}
		if !list.More {
			break
		}
		startAfter = list.Items[len(list.Items)-1].Name
	}

	if noBuckets {
		fmt.Println("No buckets")
	}

	return nil
}

func listFiles(ctx context.Context, bucket *libuplink.Bucket, prefix fpath.FPath, prependBucket bool) error {
	startAfter := ""

	for {
		bucket.ListOb
		list, err := bucket.ListObjects(ctx, storj.ListOptions{
			Direction: storj.After,
			Cursor:    startAfter,
			Prefix:    prefix.Path(),
			Recursive: *recursiveFlag,
		})
		if err != nil {
			return err
		}

		for _, object := range list.Items {
			path := object.Path
			if prependBucket {
				path = fmt.Sprintf("%s/%s", prefix.Bucket(), path)
			}
			if object.IsPrefix {
				fmt.Println("PRE", path)
			} else {
				fmt.Printf("%v %v %12v %v\n", "OBJ", formatTime(object.Modified), object.Size, path)
			}
		}

		if !list.More {
			break
		}

		startAfter = list.Items[len(list.Items)-1].Path
	}

	return nil
}

func formatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}
