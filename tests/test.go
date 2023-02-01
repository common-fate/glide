// package main

// import (
// 	"context"
// 	"fmt"

// 	cachesvc "github.com/common-fate/common-fate/pkg/cachesync"
// 	"github.com/common-fate/common-fate/pkg/deploy"
// 	"github.com/common-fate/ddb"
// )

// func run() error {
// 	ctx := context.Background()
// 	do, err := deploy.LoadConfig(deploy.DefaultFilename)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = do.LoadOutput(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	db, err := ddb.New(ctx, "common-fate-josh-7")
// 	if err != nil {
// 		return err
// 	}

// 	sync := cachesvc.CacheSyncer{
// 		DB: db,
// 	}

// 	fmt.Println("called")

// 	jsonPath := "/Users/eddie/dev/commonfate/commonfate/pkg/cachesync/schema.json"
// 	err = sync.SyncCommunityProviderSchemasFromJSON(ctx, jsonPath)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("sync complete", sync)

// 	return nil
// }

// func main() {
// 	err := run()
// 	if err != nil {
// 		fmt.Println("The error is", err)
// 	}

// }
