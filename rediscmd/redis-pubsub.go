package rediscmd

import (
	"errors"

	"github.com/go-redis/redis"
)

// Service is a service containing a pool and a connexion
type Service struct {
	client *redis.Client // conn is a redis Pool of connexions
	pubSub *redis.PubSub // pubsub interface
}

// ConnectRedis to server and create a service
func ConnectRedis(url string) (*Service, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.New("Connect: Can't parse redis url")
	}

	connexion := redis.NewClient(opt)

	if err := connexion.Ping().Err(); err != nil {
		return nil, errors.New("Connect: Can't ping redis")
	}

	return &Service{
		client: connexion,
		pubSub: nil,
	}, nil
}

// subscribe to channels
func (service *Service) subscribe(channels ...string) (<-chan *redis.Message, error) {

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
func (service *Service) SubAndManage(onMsg func(string, string) error, channels ...string) error {
	ch, err := service.subscribe(channels...)
	if err != nil {
		return err
	}
	defer service.Unsubscribe(channels...)

	for err == nil {
		select {
		case msg := <-ch:
			err = onMsg(msg.Channel, msg.Payload)
		}
	}

	return err
}

// Unsubscribe from channels
func (service *Service) Unsubscribe(channels ...string) error {
	if service.pubSub != nil {
		//defer service.pubSub.Close()
		return service.pubSub.Unsubscribe(channels...)
	}
	return errors.New("There is no Subscription to unsubscribe from")
}

// Publish publish key value to a PubSub channel
func (service *Service) Publish(channel string, value string) error {
	return service.client.Publish(channel, value).Err()
}
