package grpc

import (
"fmt"

"google.golang.org/grpc"
"google.golang.org/grpc/credentials/insecure"
dealv1 "realty/api/gen/deal/v1"
listingv1 "realty/api/gen/listing/v1"
notificationv1 "realty/api/gen/notification/v1"
searchv1 "realty/api/gen/search/v1"
userv1 "realty/api/gen/user/v1"
)

type Clients struct {
User         userv1.UserServiceClient
Listing      listingv1.ListingServiceClient
Deal         dealv1.DealServiceClient
Search       searchv1.SearchServiceClient
Notification notificationv1.NotificationServiceClient
conns        []*grpc.ClientConn
}

func NewClients(
userAddr, listingAddr, dealAddr, searchAddr, notificationAddr string,
) (*Clients, error) {
userConn, err := dial(userAddr)
if err != nil {
return nil, fmt.Errorf("подключение к user-service: %w", err)
}

listingConn, err := dial(listingAddr)
if err != nil {
return nil, fmt.Errorf("подключение к listing-service: %w", err)
}

dealConn, err := dial(dealAddr)
if err != nil {
return nil, fmt.Errorf("подключение к deal-service: %w", err)
}

searchConn, err := dial(searchAddr)
if err != nil {
return nil, fmt.Errorf("подключение к search-service: %w", err)
}

notificationConn, err := dial(notificationAddr)
if err != nil {
return nil, fmt.Errorf("подключение к notification-service: %w", err)
}

return &Clients{
User:         userv1.NewUserServiceClient(userConn),
Listing:      listingv1.NewListingServiceClient(listingConn),
Deal:         dealv1.NewDealServiceClient(dealConn),
Search:       searchv1.NewSearchServiceClient(searchConn),
Notification: notificationv1.NewNotificationServiceClient(notificationConn),
conns:        []*grpc.ClientConn{userConn, listingConn, dealConn, searchConn, notificationConn},
}, nil
}

func (c *Clients) Close() {
for _, conn := range c.conns {
conn.Close()
}
}

func dial(addr string) (*grpc.ClientConn, error) {
return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
