// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package datarepair_test

import (
	"fmt"
	"math/rand"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"storj.io/storj/internal/memory"
	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testplanet"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/uplink"
)

func TestDataRepair(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 10,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		// first, upload some remote data
		ul := planet.Uplinks[0]
		satellite := planet.Satellites[0]
		// stop discovery service so that we do not get a race condition when we delete nodes from overlay cache
		satellite.Discovery.Service.Discovery.Stop()

		satellite.Repair.Checker.Loop.Pause()
		satellite.Repair.Repairer.Loop.Pause()

		testData := make([]byte, 1*memory.MiB)
		_, err := rand.Read(testData)
		assert.NoError(t, err)

		err = ul.UploadWithConfig(ctx, satellite, &uplink.RSConfig{
			MinThreshold:     2,
			RepairThreshold:  3,
			SuccessThreshold: 4,
			MaxThreshold:     5,
		}, "testbucket", "test/path", testData)
		require.NoError(t, err)

		// get a remote segment from pointerdb
		pdb := satellite.Metainfo.Service
		listResponse, _, err := pdb.List("", "", "", true, 0, 0)
		require.NoError(t, err)

		var path string
		var pointer *pb.Pointer
		for _, v := range listResponse {
			path = v.GetPath()
			pointer, err = pdb.Get(path)
			assert.NoError(t, err)
			if pointer.GetType() == pb.Pointer_REMOTE {
				break
			}
		}

		// calculate how many storagenodes to kill
		numStorageNodes := len(planet.StorageNodes)
		redundancy := pointer.GetRemote().GetRedundancy()
		remotePieces := pointer.GetRemote().GetRemotePieces()
		minReq := redundancy.GetMinReq()
		numPieces := len(remotePieces)
		toKill := numPieces - int(minReq)
		// we should have enough storage nodes to repair on
		assert.True(t, (numStorageNodes-toKill) >= numPieces)

		// kill nodes and track lost pieces
		var lostPieces []int32
		nodesToKill := make(map[storj.NodeID]bool)
		nodesToKeepAlive := make(map[storj.NodeID]bool)
		for i, piece := range remotePieces {
			if i >= toKill {
				nodesToKeepAlive[piece.NodeId] = true
				continue
			}
			nodesToKill[piece.NodeId] = true
			lostPieces = append(lostPieces, piece.GetPieceNum())
		}

		fmt.Println("nodes to kill")
		for i := range nodesToKill {
			fmt.Println(i.String())
		}
		fmt.Println("nodes to keep alive")
		for j := range nodesToKeepAlive {
			fmt.Println(j.String())
		}

		// t.Logf("Killing %d nodes of %d", toKill, len(planet.StorageNodes))
		for _, node := range planet.StorageNodes {
			if nodesToKill[node.ID()] {
				err = planet.StopPeer(node)
				assert.NoError(t, err)
				fmt.Println("killing:", node.ID())
				// t.Logf("Killing %s = %s", node.ID(), node.Addr())
				_, err = satellite.Overlay.Service.UpdateUptime(ctx, node.ID(), false)
				assert.NoError(t, err)
			} else {
				fmt.Println("not killing:", node.ID())
				// 	// t.Logf("Keeping %s = %s", node.ID(), node.Addr())
			}
		}

		zap.L().Debug("before Checker.Loop.Restart()")
		satellite.Repair.Checker.Loop.Restart()
		satellite.Repair.Checker.Loop.TriggerWait()
		zap.L().Debug("after Checker.Loop.Triggerwait()")
		satellite.Repair.Checker.Loop.Pause()
		zap.L().Debug("before Repairer.Loop.Restart()")
		satellite.Repair.Repairer.Loop.Restart()
		satellite.Repair.Repairer.Loop.TriggerWait()
		zap.L().Debug("after Repairer.Loop.TriggerWait()")
		satellite.Repair.Repairer.Loop.Pause()
		satellite.Repair.Repairer.Limiter.Wait()
		zap.L().Debug("after Repairer.Limiter.Wait()")

		// kill nodes kept alive to ensure repair worked
		for _, node := range planet.StorageNodes {
			if nodesToKeepAlive[node.ID()] {
				err = planet.StopPeer(node)
				assert.NoError(t, err)
				fmt.Println("also killing", node.ID().String())
				// time.Sleep(5 * time.Second)

				_, err = satellite.Overlay.Service.UpdateUptime(ctx, node.ID(), false)
				assert.NoError(t, err)
			}
		}

		// we should be able to download data without any of the original nodes
		newData, err := ul.Download(ctx, satellite, "testbucket", "test/path")
		assert.NoError(t, err)
		assert.Equal(t, newData, testData)

		// updated pointer should not contain any of the killed nodes
		pointer, err = pdb.Get(path)
		assert.NoError(t, err)

		remotePieces = pointer.GetRemote().GetRemotePieces()
		fmt.Println("nodes in the pointer")
		for _, piece := range remotePieces {
			fmt.Println("in pointer", piece.NodeId)
			assert.False(t, nodesToKill[piece.NodeId])
		}
	})
}
