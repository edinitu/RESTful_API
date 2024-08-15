package main

import "testing"

func TestPopulateBalancer(t *testing.T) {
	{
		numberOfInstances = 2
		port = 7000
		balancer := new(BalancerImpl)
		populateBalancer(balancer)
		expectedAddresses := []string{"http://127.0.0.1:7001", "http://127.0.0.1:7002"}
		for idx, _ := range expectedAddresses {
			if balancer.APIInstances[idx].GetUrl().String() != expectedAddresses[idx] {
				t.Fatalf("Expected %s, got %s at index %d",
					expectedAddresses[idx], balancer.APIInstances[idx].GetUrl().String(), idx)
			}
		}
	}
	{
		numberOfInstances = 4
		port = 7000
		balancer := new(BalancerImpl)
		populateBalancer(balancer)
		expectedAddresses := []string{"http://127.0.0.1:7001", "http://127.0.0.1:7002",
			"http://127.0.0.1:7003", "http://127.0.0.1:7004"}
		for idx, _ := range expectedAddresses {
			if balancer.APIInstances[idx].GetUrl().String() != expectedAddresses[idx] {
				t.Fatalf("Expected %s, got %s at index %d",
					expectedAddresses[idx], balancer.APIInstances[idx].GetUrl().String(), idx)
			}
		}
	}
}

func TestRotate(t *testing.T) {
	{
		numberOfInstances = 2
		port = 7000
		balancer := new(BalancerImpl)
		populateBalancer(balancer)
		balancer.current = 0
		api := balancer.Rotate()
		if api.GetUrl().String() != "http://127.0.0.1:7002" {
			t.Fatalf("Expected http://127.0.0.1:7002, got %s", api.GetUrl().String())
		}
	}
}
