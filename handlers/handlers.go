package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/base-bridging-service/models"
)

func GcpPubsubPublish(cfg config.ConfigInterface) gin.HandlerFunc {
	return func(c *gin.Context) {

		if publishData, exists := c.Get("publishData"); exists {
			if topics, err := getTopics(cfg, c.FullPath()); err != nil {
				fmt.Println(err)
			} else {
				jsondata, _ := json.Marshal(publishData)

				// Publish to pubsub
				wg := new(sync.WaitGroup)
				ctx := context.Background()
				client, err := pubsub.NewClient(ctx, cfg.GetGcpProjectID())
				if err != nil {
					fmt.Println(err)
				}
				defer client.Close()

				for i := 0; i < len(topics); i++ {
					wg.Add(1)
					go publish(ctx, client, topics[i].Name, string(jsondata), wg)
				}

				wg.Wait()
			}
		}

	}
}

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	TimeStamp int64  `json:"timestamp"`
}

func publish(ctx context.Context, client *pubsub.Client, topicName, content string, wg *sync.WaitGroup) error {
	defer wg.Done()
	topic := client.Topic(topicName)

	msgUuid := uuid.New().String()

	msg, _ := json.Marshal(&Message{
		ID:        msgUuid,
		Content:   content,
		TimeStamp: time.Now().Unix(),
	})
	res := topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	_, err := res.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}

type topic struct {
	Name string `json:"name"`
}

func getTopics(cfg config.ConfigInterface, endpointPath string) ([]*topic, error) {
	ctx, rdb, topics := context.Background(), cfg.GetRedis(), []*topic{}

	// Get topics from redis
	val, err := rdb.Get(ctx, endpointPath).Result()
	if err == redis.Nil {
		// Get topics from mysql
		etrs, _ := models.NewEndpointTopicRelsTableEngine(cfg.GetDB()).GetByEndpointPath(endpointPath)

		for i := 0; i < len(etrs); i++ {
			topics = append(topics, &topic{
				Name: etrs[i].Name,
			})
		}

		jsondata, _ := json.Marshal(topics)

		// Set to redis
		if err := rdb.Set(ctx, endpointPath, string(jsondata), 3600).Err(); err != nil {
			return nil, err
		}

		return topics, nil
	} else if err != nil {
		return nil, err
	} else {
		json.Unmarshal([]byte(val), &topics)
		return topics, nil
	}
}
