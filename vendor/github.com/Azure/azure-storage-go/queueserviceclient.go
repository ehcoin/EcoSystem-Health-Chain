// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package storage

// QueueServiceClient contains operations for Microsoft Azure Queue Storage
// Service.
type QueueServiceClient struct {
	client Client
	auth   authentication
}

// GetServiceProperties gets the properties of your storage account's queue service.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/get-queue-service-properties
func (c *QueueServiceClient) GetServiceProperties() (*ServiceProperties, error) {
	return c.client.getServiceProperties(queueServiceName, c.auth)
}

// SetServiceProperties sets the properties of your storage account's queue service.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/set-queue-service-properties
func (c *QueueServiceClient) SetServiceProperties(props ServiceProperties) error {
	return c.client.setServiceProperties(props, queueServiceName, c.auth)
}
