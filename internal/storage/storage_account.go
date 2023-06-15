package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/go-autorest/autorest/to"
)

type Storage interface {
	QueryEntity(partitionKey string, rowKey string) ([]byte, error)
}

type StorageAccount struct {
	service_client *aztables.ServiceClient
}

const (
	TableName = "ShortLinks"
)

func (sa *StorageAccount) QueryEntity(partitionKey string, rowKey string) ([]byte, error) {
	client := sa.service_client.NewClient(TableName)

	filter := fmt.Sprintf("PartitionKey eq '%s' and RowKey eq '%s'", partitionKey, rowKey)
	query := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Select: to.StringPtr("RowKey,Slug,FullUrl"),
		Top: to.Int32Ptr(1),
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

func NewStorageAccount() (*StorageAccount, error) {
	sa, err := aztables.NewServiceClientFromConnectionString(os.Getenv("API_AzureStorageConnectionString"), nil)

	if err != nil {
		return nil, err
	}

	pager_options := &aztables.ListTablesOptions{
		Filter: to.StringPtr("TableName eq '" + TableName + "'"),
		Top: to.Int32Ptr(1),
	}

	pager := sa.NewListTablesPager(pager_options)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())

		if err != nil {
			return nil, err
		}

		if len(resp.Tables) == 0 {
			_, err := sa.CreateTable(context.Background(), TableName, nil)

			if err != nil {
				return nil, err
			}
		}
	}

	return &StorageAccount{service_client: sa}, err
}