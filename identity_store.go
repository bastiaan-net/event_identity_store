package store

   import (
       "context"
       "errors"
       "time"

       "go.etcd.io/etcd/clientv3"
   )

   type Store interface {
       Get(key string) (string, error)
   }

   type EtcdStore struct {
       client clientv3.Client
   }

   func NewEtcdStore(endpoints []string) (EtcdStore, error) {
       cli, err := clientv3.New(clientv3.Config{
           Endpoints:   endpoints,
           DialTimeout: 5 * time.Second,
       })
       if err != nil {
           return nil, err
       }

       return &EtcdStore{
           client: cli,
       }, nil
   }

   func (e EtcdStore) Get(key string) (string, error) {
       if err := validateKey(key); err != nil {
           return "", err
       }

       ctx, cancel := context.WithTimeout(context.Background(), 5time.Second)
       defer cancel()

       resp, err := e.client.Get(ctx, key)
       if err != nil {
           return "", err
       }

       if len(resp.Kvs) == 0 {
           return "", errors.New("key not found")
       }

       return string(resp.Kvs[0].Value), nil
   }

   func validateKey(key string) error {
       if key == "" {
           return errors.New("key cannot be empty")
       }
       return nil
   }

