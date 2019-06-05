package rediscmd

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// Service is a service containing a pool and a connexion
type service struct {
	client *redis.Client // conn is a redis Pool of connexions
	pubSub *redis.PubSub // pubsub interface
}

// ConnectRedis to server and create a service
func connectRedis(url string) (*service, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.New("Connect: Can't parse redis url")
	}

	connexion := redis.NewClient(opt)

	if err := connexion.Ping().Err(); err != nil {
		return nil, errors.New("Connect: Can't ping redis")
	}

	return &service{
		client: connexion,
		pubSub: nil,
	}, nil
}

// subscribe to channels
func (service *service) subscribe(channels ...string) (<-chan *redis.Message, error) {

	service.pubSub = service.client.Subscribe(channels...)

	// force the connection to wait, by calling the Receive() method
	iface, err := service.pubSub.Receive()
	if err != nil {
		// handle error
		return nil, errors.New("subscribe: didn't receive anything")
	}

	// Should be *Subscription, but others are possible if other actions have been
	// taken on service.pubSub since it was created.
	switch iface.(type) {
	case *redis.Subscription:
		// subscribe succeeded
		return service.pubSub.Channel(), nil
	default:
		// handle error
		return nil, errors.New("subscribe: didn't receive subscription confirmation")
	}
}

// SubAndManage subscribe and manage incoming messages from subscribed channels
// Returns when onMsg is not nil (returns an error)
func (service *service) subAndManage(onMsg func(string, string) error, channels ...string) error {
	ch, err := service.subscribe(channels...)
	if err != nil {
		return err
	}
	defer service.unsubscribe(channels...)

	for err == nil {
		select {
		case msg := <-ch:
			err = onMsg(msg.Channel, msg.Payload)
		}
	}

	return err
}

// Unsubscribe from channels
func (service *service) unsubscribe(channels ...string) error {
	if service.pubSub != nil {
		//defer service.pubSub.Close()
		return service.pubSub.Unsubscribe(channels...)
	}
	return errors.New("There is no Subscription to unsubscribe from")
}

// Publish publish key value to a PubSub channel
func (service *service) publish(channel string, value string) error {
	return service.client.Publish(channel, value).Err()
}

// SetKeyValue retrieve a Key from Redis
func (service *service) setKeyValue(key string, value interface{}, howLong time.Duration) error {
	return service.client.Set(key, value, howLong).Err()
}

// GetKey retrieve a Key from Redis
func (service *service) getKey(key string) (string, error) {
	val, err := service.client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// IncrKey incrise a key in Redis
func (service *service) incrKey(key string) (int64, error) {
	return service.client.Incr(key).Result()
}

// DecrKey incrise a key in Redis
func (service *service) decrKey(key string) (int64, error) {
	return service.client.Decr(key).Result()
}
