#! /bin/bash
#CREATES A KUBERNETES CLUSTER CALLED WP-VUL
kind create cluster --name banksim --config=kind-cluster.yml

#IMPORT CONTAINERS THAT ARE LOCAL - NEEDED for APPLE SILICONE, CAN BE REMOVED FOR AMD64
# kind load docker-image --name jaixkash876/ginauth