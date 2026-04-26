#!/bin/bash

# Prompt Queue Item
tpm-morphia-cli gen entity --schema-file ./prompt-queue/schema.yml --name BidEtPair --out-dir ../.. --with-format
tpm-morphia-cli gen entity --schema-file ./prompt-queue/schema.yml --name BucketPathPair --out-dir ../.. --with-format
tpm-morphia-cli gen entity --schema-file ./prompt-queue/schema.yml --name PromptQueueItem --out-dir ../.. --with-format
