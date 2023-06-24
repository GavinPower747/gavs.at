package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/go-autorest/autorest/to"
)

type Account interface {
	QueryEntity(partitionKey string, rowKey string) ([]byte, error)
	UpsertEntity(entity interface{}) error
}

type storageAccount struct {
	serviceClient *aztables.ServiceClient
}

const (
	TableName = "ShortLinks"
)

func (sa *storageAccount) QueryEntity(partitionKey, rowKey string) ([]byte, error) {
	client := sa.serviceClient.NewClient(TableName)

	filter := fmt.Sprintf("PartitionKey eq '%s' and RowKey eq '%s'", partitionKey, rowKey)
	query := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Select: to.StringPtr("RowKey,Slug,FullUrl"),
		Top:    to.Int32Ptr(1),
	}

	pager := client.NewListEntitiesPager(query)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())

		if err != nil {
			return []byte{}, err
		}

		for _, entity := range resp.Entities {
			return entity, nil
		}
	}

	return nil, nil
}

func (sa *storageAccount) UpsertEntity(entity interface{}) error {
	tableClient := sa.serviceClient.NewClient(TableName)

	jsonEntity, err := json.Marshal(entity)

	if err != nil {
		return err
	}

	_, err = tableClient.UpsertEntity(context.Background(), jsonEntity, nil)

	return err
}

func NewStorageAccount() (Account, error) {
	connectionString, found := os.LookupEnv("API_AzureStorageConnectionString")

	if !found {
		return nil, errors.New("API_AzureStorageConnectionString not found")
	}

	sa, err := aztables.NewServiceClientFromConnectionString(connectionString, nil)

	if err != nil {
		return nil, err
	}

	pagerOptions := &aztables.ListTablesOptions{
		Filter: to.StringPtr("TableName eq '" + TableName + "'"),
		Top:    to.Int32Ptr(1),
	}

	pager := sa.NewListTablesPager(pagerOptions)

	for pager.More() {
		resp, pageErr := pager.NextPage(context.Background())

		if pageErr != nil {
			return nil, pageErr
		}

		if len(resp.Tables) == 0 {
			_, err = sa.CreateTable(context.Background(), TableName, nil)

			if err != nil {
				return nil, err
			}
		}
	}

	return &storageAccount{serviceClient: sa}, err
}
