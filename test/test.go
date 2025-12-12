package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/go-redis/redis/v8"
)

var CachePlotBindingBitKey = "plot_binding_bitmap"

var cache = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {
	fmt.Println(SetBitMapStatus(context.Background(), 1531066177))
	fmt.Println(GetBitMapStatus(context.Background(), 1431099999))
	fmt.Println(GetBitMapStatus(context.Background(), 1599123599))
}

func GetBitMapStatus(ctx context.Context, robotID int64) (int, error) {
	// 从 Redis 读字符串
	str, err := cache.Get(ctx, CachePlotBindingBitKey).Result()
	if err != nil {
		return 0, err
	}
	if str == "" {
		return 0, nil
	}

	// Base64 解码成 []byte
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return 0, err
	}

	bm := roaring.NewBitmap()
	if err := bm.UnmarshalBinary(data); err != nil {
		return 0, err
	}

	if bm.Contains(uint32(robotID)) {
		return 1, nil
	}
	return 0, nil
}

func SetBitMapStatus(ctx context.Context, robotID int64) error {
	bm := roaring.NewBitmap()

	// 从 Redis 读取已有数据（字符串）
	oldStr, _ := cache.Get(ctx, CachePlotBindingBitKey).Result()
	if oldStr != "" {
		oldData, _ := base64.StdEncoding.DecodeString(oldStr)
		_ = bm.UnmarshalBinary(oldData)
	}

	// 添加新的 ID
	bm.Add(uint32(robotID))

	// 序列化并转成 Base64 字符串
	raw, _ := bm.ToBytes()
	encoded := base64.StdEncoding.EncodeToString(raw)

	// 存入 Redis（string）
	err := cache.Set(ctx, CachePlotBindingBitKey, encoded, 0).Err()
	return err
}
